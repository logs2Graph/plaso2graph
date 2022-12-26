# Roadmap

This is the roadmap of the project, it describe the current state of the project and the next steps.
It shows the differents Artefacts and Logs handle by the project and the one that are not yet implemented.

It also shows the differents relationships between the entities and the ones that will be implemented.

- Process:
  - [x] Evtx EventID 4688
  - [x] Evtx Sysmon EventID 1
  - [x] Prefetch
  - [x] UserAssist
  - [x] ShellBag
  - [x] SRUM
  - [ ] MRU
  - [x] AMCache
  - [ ] ShimCache
  - [x] AppCompatCache
  - [x] Evtx Sysmon EventID 10 (memory access)
  - [x] Evtx Sysmon EventID 7 (Image Loaded)
  - [x] Evtx Sysmon EventID 9 (Raw Access Read)
  - [x] Lnk (shortcut)
- Scripts:
  - [x] Evtx EventID 4103 // TODO: Parse ContextInfo
  - [x] Evtx EventID 4104
  - [ ] Powershell Transcript
- User:
  - [x] Evtx Security
  - [x] Evtx Sysmon
  - [ ] Evtx EventID 4673 (Log User's Privileges)
  - [x] SAM User Registry
- Groups:
  - [ ] Evtx EventID 4627 (Group Membership)
- Login Events:
  - [x] Evtx EventID 4648 (Explicit Credentials)
  - [ ] Evtx EventID 4649 (Replay Attack Detected)
  - [x] Evtx EventID 4624 (User Logon)
  - [x] Evtx EventID 4625 (User Fail to Logon)
  - [x] Evtx EventID 4634 (User Logoff)
- Computer:
  - [x] Evtx
- Connection:
  - [x] Evtx Sysmon EventID 3
  - [ ] Evtx EventID 5031
  - [ ] SRUM Connectivity?
- WebHistory
  - [x] Chrome
  - [x] Firefox
- AutoRun:
  - Scheduled Tasks
    - [x] Windows Jobs
    - [ ] Evtx EventID 4699-4702
  - [x] Registry Run / RunOnce
  - [ ] Boot Execution
  - [x] Service
  - [ ] Task Cache
- File:
  - [x] Evtx Sysmon Event ID 9 (Raw Access Read) (Not Tested)
  - [x] Evtx Sysmon Event ID 11 (File Create) 
  - [x] Evtx Sysmon Event ID 23 (File Delete)
  - [x] MFT (Note: Works but will take a long time.)
- Misc Events:
  - [ ] USB
- Changes:
  - AV Disabled
  - FW Changes
  - GPO Changes
  - User Changes
    - [x] EventID 4738 (A user account was changed)
    - [x] EventID 4720 (User Created)
    - [x] EventID 4722 (User Enabled)
    - [x] EventID 4723-4724 (Password reset or changed)
    - [x] EventID 4725 (User Disabled)
    - [x] EventID 4726 (User Deleted)
    - [ ] EventID 4704 (User right assigned)
    - [ ] EventID 4705 (User right removed)
  - Group Changes
    - [x] EventID	4727 	A security-enabled global group was created
    - [x] EventID	4728 	(member was added to a security-enabled global group)
    - [x] EventID	4729 	(member was removed from a security-enabled global group)
    - [x] EventID	4730 	(security-enabled global group was deleted)
    - [x] EventID	4731 	(security-enabled local group was created)
    - [x] EventID	4732 	(member was added to a security-enabled local group)
    - [x] EventID	4733 	(member was removed from a security-enabled local group)
    - [x] EventID	4734 	(security-enabled local group was deleted)
    - [x] EventID	4735 	(security-enabled local group was changed)
    - [x] EventID	4737 	(security-enabled global group was changed)
- Exporter
  - [x] Neo4j
  - [] Json
  - [] Xml
  - [] Csv
- Improvements
- [x] Convert Event Entities to Relationships
- [] Merge Files when possible
- [] Assign Computer Name by default (To support multiple computers)


Optimisation Ideas:
- Perform Merge inside "Add" Functions

Relationships:
- [x] User -[DELETE]->User
- [x] User -[DELETE]->Group
- [x] User -[CREATE]->User
- [x] User -[CREATE]->Group
- [x] User -[CHANGE]->User
- [x] User -[CHANGE]->Group
- [x] User -[ENABLE]->User
- [x] User -[ENABLE]->Group
- [x] User -[DISABLE]->User
- [x] User -[DISABLE]->Group
- [x] User -[LOGON]->Computer
- [x] User -[LOGOFF]->Computer
- [x] User -[LOGON]->User
- [ ] User -[ACCESS]->File
