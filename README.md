# aeaconf2

new config format for [aeacus](https://github.com/elysium-suite/aeacus)

## syntax

- the beginning part is ini
- the check/condition definitions use the custom language (see comments)

```
[round]
title = Linux ICC
os = Ubuntu 20.04.03
user = cpadmin
local = false

[remote]
enable = true
name = LinICC
server = https://scoring.cyberaegis.tech
password = password
---
// check messages must be enclosed in quotes
"Check1: cool vuln!": 2
  // condition lines MUST be indented

  // functions may have individual hints, enclosed in square brackets
  ServiceUp "sshd" ["Install SystemD!"]

  // lines that don't end in a boolean operator are implicitly ANDed
  PathExists "/etc/ssh"

// entire checks may have hints
"Check2": -3 ["Is the system broken?"]
  // lines that end in boolean operators continue the expression on the next line
  ServiceUp "sshd" && PathExistsNot "/abc" || ServiceUp "samba" ||

  // parentheses may be used as you would in any other boolean expression
  (PathExists "/etc/ssh" || PathExists "/var") && ServiceUp "abc"

// checks can be shortened to a single line with a semicolon (;)
"Check3": 4 ["hint before semicolon"]; ServiceUp "sshd" || PathExists "/etc/ssh"

"Check4": 6
  // hints can also be specified on parenthesized expressions
  (ServiceUp "sshd" && PathExists "/") ["Did you delete /?"]
```

# parsed

this is the parsed output (it parses hints as well on a per-check, per-condition, and even per-boolean-expression level)

```
Check1: cool vuln! : 2
  AND {
    ServiceUp(Service="sshd")
    PathExists(Path="/etc/ssh")
  }

Check2 : -3
  OR {
    OR {
      AND {
        ServiceUp(Service="sshd")
        NOT {
          PathExists(Path="/abc")
        }
      }
      ServiceUp(Service="samba")
    }
    AND {
      OR {
        PathExists(Path="/etc/ssh")
        PathExists(Path="/var")
      }
      ServiceUp(Service="abc")
    }
  }

Check3 : 4
  OR {
    ServiceUp(Service="sshd")
    PathExists(Path="/etc/ssh")
  }

Check4 : 6
  AND {
    ServiceUp(Service="sshd")
    PathExists(Path="/")
  }
```
