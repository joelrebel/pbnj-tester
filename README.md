#### PBnJ tester

The pbnj-tester can be used to test various PBnJ functionality on one or servers BMCs.

The plan is to have the pbnj-tester perform an action like `power-cycle` and
validate that the PBnJ action was successful and that the server was power cycled.

Note: status beta, there be bugs.

### list supported tests

```
❯ ./pbnj-tester list-tests
power-status
power-on
power-off
power-cycle
pxeBoot
```

### run tests on a single server

```json
❯ ./pbnj-tester run --pbnj-addr localhost:50051 --bmc-ip 192.168.1.1 --bmc-user root --bmc-pass hunter2 --tests power-status,power-on
[
 {
  "TestName": "power-status",
  "Output": null,
  "Error": null,
  "Succeeded": true,
  "Runtime": 0
 },
 {
  "TestName": "power-on",
  "Output": null,
  "Error": null,
  "Succeeded": true,
  "Runtime": 0
 }
]
```

### run tests on multiple servers

This command requires a test configuration, see the [tests.yaml](tests.yaml) for a sample test configuration.

The tests are executed concurrently and a JSON output is returned.

```json
❯ ./pbnj-tester run-multiple --pbnj-addr localhost:50051 --tests-config tests.yaml
[
 {
  "Vendor": "ASRockRack",
  "Model": "ROMED8HM3 EPYC 7402P",
  "Name": "foo",
  "BMCIP": "192.168.1.1",
  "Results": [
   {
    "TestName": "power-status",
    "Output": null,
    "Error": null,
    "Succeeded": true,
    "Runtime": 0
   }
  ]
 },
 {
  "Vendor": "Super Micro",
  "Model": "SSG-110P-NTR10",
  "Name": "bar",
  "BMCIP": "192.168.1.2",
  "Results": [
   {
    "TestName": "power-status",
    "Output": null,
    "Error": null,
    "Succeeded": true,
    "Runtime": 0
   }
  ]
 }
]
```


### TODO

- Results to include the method used by the test when it was successful, the
method being Redfish/IPMI/SSH/Vendor API.
