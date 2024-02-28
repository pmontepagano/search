package contract

import (
	"bytes"
	"crypto/sha512"
	"fmt"
	"sync"

	"github.com/pmontepagano/search/cfsm"
	"github.com/vishalkuo/bimap"

	pb "github.com/pmontepagano/search/gen/go/search/v1"
)

type ContractOutputFormats int

const (
	SingleCFSMPythonBisimulation ContractOutputFormats = iota
)

type Contract interface {
	GetContractID() string
	GetRemoteParticipantNames() []string // Returns the names of all participants in this contract (except the Service Provider, who is unnamed).
	GetBytesRepr() []byte
	Convert(ContractOutputFormats) ([]byte, error)
}

// LocalContract is an interface that represents a local view of a contract. It is used to
// specify behaviour for Service Providers in SEARch.
type LocalContract interface {
	Contract
	GetFormat() pb.LocalContractFormat
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

type GlobalCFSMContract struct {
	*cfsm.System
	id               string
	localParticipant *cfsm.CFSM
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

func (lc *LocalCFSMContract) GetRemoteParticipantNames() []string {
	return lc.CFSM.OtherCFSMs()
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

func (c *LocalCFSMContract) GetBytesRepr() []byte {
	return c.Bytes()
}

func (c *GlobalCFSMContract) GetFormat() pb.GlobalContractFormat {
	return pb.GlobalContractFormat_GLOBAL_CONTRACT_FORMAT_FSA
}

func (lc *LocalCFSMContract) GetFormat() pb.LocalContractFormat {
	return pb.LocalContractFormat_LOCAL_CONTRACT_FORMAT_FSA
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

func (lc *LocalCFSMContract) Convert(format ContractOutputFormats) ([]byte, error) {
	if format == SingleCFSMPythonBisimulation {
		code, _, _, err := cfsm.ConvertCFSMToPythonBisimulationFormat(lc.CFSM)
		return code, err
	}
	return nil, fmt.Errorf("invalid output format for this type of contract")
}

func (lc *GlobalCFSMContract) Convert(format ContractOutputFormats) ([]byte, error) {
	return nil, fmt.Errorf("invalid output format for this type of contract")
}

type LocalPyCFSMContract struct {
	pythonCode []byte
	id         string
	sync.Mutex
	convertedFrom               *LocalCFSMContract
	participantNameTranslations *bimap.BiMap[string, string]
	messageTranslations         *bimap.BiMap[string, string]
}

func (lc *LocalCFSMContract) ConvertToPyCFSM() (*LocalPyCFSMContract, error) {
	pythonCode, participantTranslations, messageTranslations, err := cfsm.ConvertCFSMToPythonBisimulationFormat(lc.CFSM)
	if err != nil {
		return nil, err
	}
	return &LocalPyCFSMContract{
		pythonCode:                  pythonCode,
		convertedFrom:               lc,
		participantNameTranslations: participantTranslations,
		messageTranslations:         messageTranslations,
	}, nil
}
