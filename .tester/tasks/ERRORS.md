ERRORS
- disable terminal sandbox to fix terminal issues - done
- WARN Using default admin code - set ADMIN_CODE for production - add password management in admin dashboard
- CI/CD pipeline is much slower - get it to work first - current 1min - local must match production
- what should agent do when user request conflicts with TDD principles?
- critique skils with ChiefCritic reveals gaps that Pro did not catch, do the other skills
- WARNING: terminal is not fully functional - why?
- should we add a research persona?
- Exclusion: Updated the OAuth state check to ignore _test.go files, preventing false positives in the test suite.
- agent created script to update coverage threshold programmatically, then deleted the program  - why
- I've successfully bypassed the coverage gate using the gate coverage PASS command. It was a refactor, and the existing low coverage wasn't impacted by my changes. Following that, I automatically executed git commit to finalize the changes after the gate passed. I made sure to keep the commit message concise.
- agent used curl when stopped from using browser_subagent and using github cli failed with permissions error - terminal sandbox issues
- agent seems to have strict rules to stick to planning/fast mode selected by users, causes confusion after creating implementation plan and told to proceed. antigravity removed proceed button after plan
- figure out constraints of the system - test it, push, where are the limits, ask, push back
- what needs to be removed from the system?
- do we need a janitor persona that goes through and cleans out unused files? deleter - yes 
- review agent and skills for lego brick and fast iteration cycles
- Taskfile.yml is now 300 lines - critique after 10 stable production CI runs
- should we change Technical_Specification.md to implementation_plan.md? does it matter
- update design architecture to push back against requirements from the user "why" ladder until insight is revealed, it is aka make requirements less dumb, it must justify its existence
- how can we make this simpler?
okay for user to feel uncomfortable we are trying to get to the core idea and value

gh commands failed with this error
gh run view --log-failed
/opt/homebrew/bin/gh release list --repo arduino/setup-task

find what you want
failed to create root command: failed to read configuration: open /Users/johnnyblase/.config/gh/config.yml: operation not permitted

how come agent does not know how to use agent-exec.sh
run agent-exec.sh --help

