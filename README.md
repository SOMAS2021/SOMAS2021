# SOMAS2021

## Set up

1. Download Go 1.17 from [here](https://go.dev/doc/install).
2. Download `golangci-lint` from [here](https://golangci-lint.run/usage/install/).

## Contribution Guidelines
A lot of these guidelines are from the [SOMAS2020 repo](https://github.com/SOMAS2020/SOMAS2020/blob/main/docs/SETUP.md) :)
### Coding Rules
1. You're encouraged to use [VSCode](https://code.visualstudio.com/) with the [Go extension](https://code.visualstudio.com/docs/languages/go).
2. Trust the language server. Red lines == death. Yellow lines == close to death. An example where it might be very tempting to let yellow lines pass are in `struct`s:
```golang
type S struct {
    name string
    age  int
}

s1 := S {"pitt", 42} // NO
s2 := S {"pitt"} // NO (even though it initialises age to 0)

s3 := S{name: "pitt", age: 42} // OK, if we add fields into S, this will still be correct
s4 := S{name: "pittson"} // OK if `pittson`'s age is 0
```
3. Write tests where required. There are many guides online on how to do this in Golang. Tests will be run alongside CI when you pull into the main repo. If anyone breaks anything, it's easy to observe that if you have tests. Otherwise, your code will be broken unknowingly.

4. DO NOT TOUCH code you don't own unless you have a good reason to. If you have a good reason to, do it in a separate PR and notify the owners of the code.

5. Do not use `panic` or `die` - return an `error` instead!

6. Do not use system-specific packages (e.g. `internal/syscall/unix`).

7. Use the superior `errors.Errorf` to create your errors so that we have a stack trace.
### Code Reviews and PRs
- Do not push to the `main` branch. Make your own branch and create a PR into `main` when it's ready for review.
- When working on your own team's features, please name your branch as: `teamX/FEATURE_NAME/WHATEVER_YOU_LIKE_HERE`
- When ready to merge into your team's feature branch, create a PR to merge into `teamX/FEATURE_NAME`. When the feature is complete, then create a PR into `main`.
- Keep PRs small: if they are too large, you will be told to split your code into smaller PRs. This ensures they can be reviewed properly.
- You need to make sure your code is up-to-date with the `main` branch. Merge commits are *not allowed*: learn how to rebase. Ask someone in the #git-env-infra channel if you don't know how to do this.
- Do not review your own code. That completely defeats the purpose of code review.
- Team leads: when doing a PR for your team's code, try not review it yourself - get another team lead to review it.
- Review PRs in a timely manner! Ideally by the next day so that other teams aren't blocked.
## Team Leads
1. Jaafar Rammal (JaafarRammal) and Cyrus Goodarzi (Silvertrousers)
2. Tom Eaton (tomeaton17)
3. Sara Fernandez (SaraFFdez)
4. Hussain Kurabadwala (hussain2603)
5. Jason Zheng (jzzheng22)
6. Matt Scott (MattScottEEE)
7. Moin Bukhari (moinbukhari)

