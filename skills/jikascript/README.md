# JikaScript Agent Skill

A skills.sh-compatible agent skill for the JikaScript language, compiler, CLI,
and cooperative scheduler. It is a deep reference for the language constructs
(the main focus) with exhaustive, example-rich material, plus install, build,
and flash guidance.

## Install

The skill lives in this repo so it can be cloned and installed on any machine.
The standard skills CLI installs it from this repo (add `-g` for a global
install across projects; `-y` skips prompts):

```bash
npx skills add sayanriju/jikascript --skill jikascript
# global, non-interactive:
npx skills add sayanriju/jikascript --skill jikascript -g -y
```

`npx skills add` handles the agent's skills directory, project vs global scope,
and updates. Your agent loads the skill by the `name` and `description` in the
`SKILL.md` frontmatter.

### Manual fallback (no Node/npm)

Install into your agent/skills.sh skills directory (typically
`~/.agents/skills`), either by copying or symlinking:

```bash
cd jikascript                                  # repo root
mkdir -p ~/.agents/skills
ln -s "$PWD/skills/jikascript" ~/.agents/skills/jikascript
```

Or copy it so edits are independent of the clone:

```bash
cp -r skills/jikascript ~/.agents/skills/jikascript
```

Manual installs do not participate in `npx skills update`. On a fresh clone,
re-run the link line (or `cp -r`) to keep the installed copy present after a
`git clone`.

Uninstall:

```bash
rm ~/.agents/skills/jikascript        # or rm -rf for a copied install
# or via the CLI:
npx skills remove jikascript
```

## What's inside

```
skills/jikascript/
├── SKILL.md                 # skill metadata + usage guidance + index
├── examples/                # copy-ready .jika programs
│   ├── blink.jika
│   ├── multiple_tasks.jika
│   ├── conditional_task.jika
│   ├── objects.jika
│   └── multi-file/{app,helpers}.jika
└── references/
    ├── README.md            # reference index + core facts
    ├── language.md          # exhaustive language construct reference
    ├── arrays.md            # bounded arrays
    ├── objects.md           # static object records -> C structs
    ├── tasks.md             # cooperative periodic task model
    ├── scheduler.md         # standalone pure-C scheduler API
    ├── hardware.md          # GPIO hardware API
    ├── includes.md          # multi-file include splicing
    ├── cli.md               # install, build, transpile, check, flash
    └── example-programs.md  # complete end-to-end programs
```

The authoritative specification lives at the repository root in `SPEC.md`;
this skill is a practical reference built from it and the live examples.