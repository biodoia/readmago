#!/bin/bash
# eval-cli-agent.sh — Valutazioni eseguite da agenti CLI esterni con varie combo di skill e LLM
#
# Ogni combinazione (agent CLI, skill, modello LLM) produce una valutazione
# dell'ecosistema fgt-sdk (design/taste/anti-slop/readme). Le valutazioni vengono
# scritte in:
#   /home/lisergico25/projects/readmago/kb/evals/<ts>-<combo>.md
# e iniettate nel meta-dev loop come prossimo comando di miglioramento.
#
# Combos supportate (default matrix):
#   codex + gofrontai skill + gpt-5.6        (design eval)
#   agy   + readme-gen  skill + gemini-3.6-flash-low (readme gen eval)
#   grok  + readme-lint  skill + grok-4.5     (lint/eval)
#   codex + no-ai-slop writing eval + gpt-5.6

set -u

EVAL_DIR=/home/lisergico25/projects/readmago/kb/evals
mkdir -p "$EVAL_DIR"
TS=$(date +%Y%m%dT%H%M%SZ)

run_combo() {
  local cli="$1"      # codex | agy | grok
  local skill="$2"    # skill name (gofrontai | readme-gen | readme-lint | no-ai-slop | ...)
  local model="$3"    # LLM model override
  local agent="$4"    # target repo / focus
  local prompt="$5"

  local out="$EVAL_DIR/${TS}-${cli}-${skill}-${model}.md"
  {
    echo "# EVAL ${TS} — combo: cli=${cli} skill=${skill} llm=${model} target=${agent}"
    echo ""
    echo "## Input"
    echo ""
    echo "$prompt"
    echo ""
    echo "## Combo"
    echo ""
    echo "- Agent CLI: \`${cli}\`"
    echo "- Skill: \`${skill}\`"
    echo "- Modello LLM: \`${model}\`"
    echo "- Target: \`${agent}\`"
    echo ""
    echo "## Setup"
    echo ""
    echo "- skill dir: /home/lisergico25/projects/agent-skills/${skill} (o equivalente)"
    echo "- binary: /home/lisergico25/.local/bin/${cli}"
    echo "- bind: loopback; aigoproxy per esposizione"
    echo ""
  } > "$out"

  local rc=0
  case "$cli" in
    codex)  /home/lisergico25/.local/bin/codex exec --model "$model" --skill "$skill" "$prompt" >> "$out" 2>&1 || rc=$? ;;
    agy)    /home/lisergico25/.local/bin/agy --model "$model" --skill "$skill" "$prompt" >> "$out" 2>&1 || rc=$? ;;
    grok)   /home/lisergico25/.local/bin/grok --model "$model" --skill "$skill" "$prompt" >> "$out" 2>&1 || rc=$? ;;
  esac

  echo "EVAL_DONE ts=$TS cli=$cli skill=$skill llm=$model target=$agent rc=$rc report=$out"
}

# Matrix di valutazioni esterne (skill × LLM × target)
run_combo codex "gofrontai" "gpt-5.6-sol" "gofrontai-design" \
  "Valuta la qualità design del repo gofrontai secondo la skill gofrontai: anti-slop, archetype, wireframe, firme visive. Report JSON con score 1-10."

run_combo agy "readme-gen" "gemini-3.6-flash-low" "readmago-readme" \
  "Genera un README per il repo biodoia/readmago secondo la skill readme-gen e valuta anti-slop: output deve superare il quality gate >= 4/5."

run_combo grok "readme-lint" "grok-4.5" "readmago-readme-lint" \
  "Valuta la bontà di un README generato secondo la skill readme-lint: 7-axis quality gate, no-ai-slop, diagram self-contained HTML+SVG."

echo "=== EVAL RUN COMPLETE $(date -Iseconds) ==="
