package contract

import (
	"bytes"
	"crypto/sha512"
	"fmt"
	"sync"
	"text/template"

	"github.com/pmontepagano/search/cfsm"
	"github.com/vishalkuo/bimap"

	pb "github.com/pmontepagano/search/gen/go/search/v1"
)

type Contract interface {
	GetContractID() string
	GetRemoteParticipantNames() []string // Returns the names of all participants in this contract (except the Service Provider, who is unnamed).
	GetBytesRepr() []byte
}

// LocalContract is an interface that represents a local view of a contract. It is used to
// specify behaviour for Service Providers in SEARch.
type LocalContract interface {
	Contract
	GetFormat() pb.LocalContractFormat
	Convert(pb.LocalContractFormat) (LocalContract, error)
	// TODO: add GetNextState()
}

// GlobalContract is an interface that represents the global view of a communication channel.
// It is used to specify behaviour for Service Clients in SEARch.
// You can project a GlobalContract into one LocalContract for each participant in the GlobalContract.
type GlobalContract interface {
	Contract
	GetFormat() pb.GlobalContractFormat
	GetParticipants() []string                   // Returns the names of all participants in this contract.
	GetLocalParticipantName() string             // Returns the name of the local participant of this contract.
	GetProjection(string) (LocalContract, error) // Returns the LocalContract for the given participant name.
}

type LocalCFSMContract struct {
	*cfsm.CFSM
	id string
	sync.Mutex
}

func (c *LocalCFSMContract) GetContractID() string {
	locked := c.TryLock()
	if locked {
		defer c.Unlock()
		if c.id != "" {
			return c.id
		}
	}
	contractHash := sha512.Sum512(c.GetBytesRepr())
	id := fmt.Sprintf("%x", contractHash[:])
	if locked {
		c.id = id
	}
	return id
}

func (lc *LocalCFSMContract) GetRemoteParticipantNames() []string {
	return lc.CFSM.OtherCFSMs()
}

func (c *LocalCFSMContract) GetBytesRepr() []byte {
	return c.Bytes()
}

func (lc *LocalCFSMContract) GetFormat() pb.LocalContractFormat {
	return pb.LocalContractFormat_LOCAL_CONTRACT_FORMAT_FSA
}

func (lc *LocalCFSMContract) Convert(format pb.LocalContractFormat) (LocalContract, error) {
	if format == pb.LocalContractFormat_LOCAL_CONTRACT_FORMAT_PYTHON_BISIMULATION_CODE {
		return lc.ConvertToPyCFSM()
	}
	return nil, fmt.Errorf("invalid output format for this type of contract")
}

type GlobalCFSMContract struct {
	*cfsm.System
	id               string
	localParticipant *cfsm.CFSM
	sync.Mutex
}

func (c *GlobalCFSMContract) GetContractID() string {
	locked := c.TryLock()
	if locked {
		defer c.Unlock()
		if c.id != "" {
			return c.id
		}
	}
	contractHash := sha512.Sum512(c.GetBytesRepr())
	id := fmt.Sprintf("%x", contractHash[:])
	if locked {
		c.id = id
	}
	return id
}

func (c *GlobalCFSMContract) GetParticipants() (ret []string) {
	for _, m := range c.CFSMs {
		ret = append(ret, m.Name)
	}
	return
}

func (c *GlobalCFSMContract) GetLocalParticipantName() string {
	return c.localParticipant.Name
}

func (c *GlobalCFSMContract) GetRemoteParticipantNames() (ret []string) {
	return c.getParticipants(false)
}

func (c *GlobalCFSMContract) getParticipants(includeLocal bool) (ret []string) {
	for _, m := range c.CFSMs {
		if !includeLocal && m == c.localParticipant {
			continue
		}
		ret = append(ret, m.Name)
	}
	return
}

func (c *GlobalCFSMContract) GetBytesRepr() []byte {
	return c.Bytes()
}

func (c *GlobalCFSMContract) GetFormat() pb.GlobalContractFormat {
	return pb.GlobalContractFormat_GLOBAL_CONTRACT_FORMAT_FSA
}

func (c *GlobalCFSMContract) GetProjection(participantName string) (LocalContract, error) {
	// TODO: Do a deep copy of the CFSM and remove pointers to other machines.
	machine, err := c.GetMachine(participantName)
	if err != nil {
		return nil, err
	}
	contract := LocalCFSMContract{
		CFSM: machine,
	}
	return &contract, nil
}

func ConvertPBGlobalContract(pbContract *pb.GlobalContract) (GlobalContract, error) {
	if pbContract.Format == pb.GlobalContractFormat_GLOBAL_CONTRACT_FORMAT_FSA {
		cfsmSystem, err := cfsm.ParseSystemCFSMsFSA(bytes.NewReader(pbContract.Contract))
		if err != nil {
			return nil, err
		}
		for _, m := range cfsmSystem.CFSMs {
			if m.Name == pbContract.InitiatorName {
				contract := GlobalCFSMContract{
					System:           cfsmSystem,
					localParticipant: m,
				}
				return &contract, nil
			}
		}
		return nil, fmt.Errorf("initiator name not found in contract")
	}
	return nil, fmt.Errorf("not implemented")
}

func ConvertPBLocalContract(pbContract *pb.LocalContract) (LocalContract, error) {
	if pbContract.Format == pb.LocalContractFormat_LOCAL_CONTRACT_FORMAT_FSA {
		machine, err := cfsm.ParseSingleCFSMFSA(bytes.NewReader(pbContract.Contract))
		if err != nil {
			return nil, err
		}
		contract := LocalCFSMContract{
			CFSM: machine,
		}
		return &contract, nil
	}
	return nil, fmt.Errorf("not implemented")
}

type LocalPyCFSMContract struct {
	pythonCode []byte
	id         string
	sync.Mutex
	convertedFrom               *LocalCFSMContract
	participantNameTranslations *bimap.BiMap[string, string]
	messageTranslations         *bimap.BiMap[string, string]
	Selfname                    string
}

func (lc *LocalCFSMContract) ConvertToPyCFSM() (*LocalPyCFSMContract, error) {
	pythonCode, participantTranslations, messageTranslations, selfname, err := cfsm.ConvertCFSMToPythonBisimulationFormat(lc.CFSM)
	if err != nil {
		return nil, err
	}
	return &LocalPyCFSMContract{
		pythonCode:                  pythonCode,
		convertedFrom:               lc,
		participantNameTranslations: participantTranslations,
		messageTranslations:         messageTranslations,
		Selfname:                    selfname,
	}, nil
}

// TODO: this is copy-pasted from LocalCFSMContract's implementation. Refactor w/ a common interface.
func (c *LocalPyCFSMContract) GetContractID() string {
	locked := c.TryLock()
	if locked {
		defer c.Unlock()
		if c.id != "" {
			return c.id
		}
	}
	contractHash := sha512.Sum512(c.GetBytesRepr())
	id := fmt.Sprintf("%x", contractHash[:])
	if locked {
		c.id = id
	}
	return id
}

func (lc *LocalPyCFSMContract) GetRemoteParticipantNames() []string {
	inverseMap := lc.participantNameTranslations.GetInverseMap()
	res := make([]string, len(inverseMap))
	i := 0
	for k := range inverseMap {
		res[i] = k
		i++
	}
	return res
}

func (c *LocalPyCFSMContract) GetBytesRepr() []byte {
	return c.pythonCode
}

func (lc *LocalPyCFSMContract) GetFormat() pb.LocalContractFormat {
	return pb.LocalContractFormat_LOCAL_CONTRACT_FORMAT_PYTHON_BISIMULATION_CODE
}

func (lc *LocalPyCFSMContract) Convert(format pb.LocalContractFormat) (LocalContract, error) {
	return nil, fmt.Errorf("not implemented")
}

func (lc *LocalPyCFSMContract) GetPythonCode(varName string) string {
	type templateData struct {
		MachineName string
	}
	// Parse the template string
	tmpl, err := template.New("templatePyCFSM").Parse(string(lc.pythonCode))
	if err != nil {
		panic(err)
	}

	// Create a buffer to capture the output
	var buffer bytes.Buffer

	// Execute the template
	err = tmpl.Execute(&buffer, templateData{MachineName: varName})
	if err != nil {
		panic(err)
	}
	return buffer.String()
}

func (lc *LocalPyCFSMContract) GetOriginalParticipantName(translated string) (string, error) {
	original, ok := lc.participantNameTranslations.GetInverse(translated)
	if !ok {
		return "", fmt.Errorf("translated name not found")
	}
	return original, nil
}
