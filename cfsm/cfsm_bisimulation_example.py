from cfsm_bisimulation import CommunicatingFiniteStateMachine
from z3 import Int

cfsm_clientapp = CommunicatingFiniteStateMachine(['Clientapp', 'Srv', 'Pps'])
cfsm_clientapp.add_states('clientq0', 'q1', 'q2', 'q3', 'q4', 'q5', 'q6', 'q7')
cfsm_clientapp.set_as_initial('clientq0')
cfsm_clientapp.add_transition_between(
    'clientq0',
    'q1',
    'Clientapp Srv ! PurchaseRequest()'
)
cfsm_clientapp.add_transition_between(
    'q1',
    'q2',
    'Clientapp Srv ? TotalAmount()'
)
cfsm_clientapp.add_transition_between(
    'q2',
    'q3',
    'Clientapp Pps ! CardDetailsWithTotalAmount()'
)
cfsm_clientapp.add_transition_between(
    'q3',
    'q4',
    'Clientapp Pps ? PaymentNonce()'
)
cfsm_clientapp.add_transition_between(
    'q4',
    'q5',
    'Clientapp Srv ! PurchaseWithPaymentNonce()'
)
cfsm_clientapp.add_transition_between(
    'q5',
    'q6',
    'Clientapp Srv ? PurchaseOK()'
)
cfsm_clientapp.add_transition_between(
    'q5',
    'q7',
    'Clientapp Srv ? PurchaseFail()'
)



cfsm_srv = CommunicatingFiniteStateMachine(['Srv', 'Clientapp', 'Pps'])
cfsm_srv.add_states('srvq0', 'q1', 'q2', 'q3', 'q4', 'q5', 'q6', 'q7', 'q8')
cfsm_srv.set_as_initial('srvq0')
cfsm_srv.add_transition_between(
    'srvq0',
    'q1',
    'Srv Clientapp ? PurchaseRequest()'
)
cfsm_srv.add_transition_between(
    'q1',
    'q2',
    'Srv Clientapp ! TotalAmount()'
)
cfsm_srv.add_transition_between(
    'q2',
    'q3',
    'Srv Clientapp ? PurchaseWithPaymentNonce()'
)
cfsm_srv.add_transition_between(
    'q3',
    'q4',
    'Srv Pps ! RequestChargeWithNonce()'
)
cfsm_srv.add_transition_between(
    'q4',
    'q5',
    'Srv Pps ? ChargeOK()'
)
cfsm_srv.add_transition_between(
    'q4',
    'q6',
    'Srv Pps ? ChargeFail()'
)
cfsm_srv.add_transition_between(
    'q5',
    'q7',
    'Srv Clientapp ! PurchaseOK()'
)
cfsm_srv.add_transition_between(
    'q6',
    'q8',
    'Srv Clientapp ! PurchaseFail()'
)



cfsm_pps = CommunicatingFiniteStateMachine(['Pps', 'Clientapp', 'Srv'])
cfsm_pps.add_states('ppsq0', 'q1', 'q2', 'q3', 'q4', 'q5')
cfsm_pps.set_as_initial('ppsq0')
cfsm_pps.add_transition_between(
    'ppsq0',
    'q1',
    'Pps Clientapp ? CardDetailsWithTotalAmount()'
)
cfsm_pps.add_transition_between(
    'q1',
    'q2',
    'Pps Clientapp ! PaymentNonce()'
)
cfsm_pps.add_transition_between(
    'q2',
    'q3',
    'Pps Srv ? RequestChargeWithNonce()'
)
cfsm_pps.add_transition_between(
    'q3',
    'q4',
    'Pps Srv ! ChargeOK()'
)
cfsm_pps.add_transition_between(
    'q3',
    'q5',
    'Pps Srv ! ChargeFail()'
)


cfsm_pps_alt = CommunicatingFiniteStateMachine(['Ppsalt', 'Cliente', 'Backend', ])
cfsm_pps_alt.add_states('ppsaltq0', 'q1', 'q2', 'q3', 'q4', 'q5')
cfsm_pps_alt.set_as_initial('ppsaltq0')
cfsm_pps_alt.add_transition_between(
    'ppsaltq0',
    'q1',
    'Ppsalt Cliente ? CardDetailsWithTotalAmount()'
)
cfsm_pps_alt.add_transition_between(
    'q1',
    'q2',
    'Ppsalt Cliente ! PaymentNonce()'
)
cfsm_pps_alt.add_transition_between(
    'q2',
    'q3',
    'Ppsalt Backend ? RequestChargeWithNonce()'
)
cfsm_pps_alt.add_transition_between(
    'q3',
    'q4',
    'Ppsalt Backend ! ChargeOK()'
)
cfsm_pps_alt.add_transition_between(
    'q3',
    'q5',
    'Ppsalt Backend ! ChargeFail()'
)


relation, matches = cfsm_pps.calculate_bisimulation_with(cfsm_pps_alt)

assert bool(relation)

assert {'Pps': 'Ppsalt', 'Srv': 'Backend', 'Clientapp': 'Cliente'} == matches['participants']
