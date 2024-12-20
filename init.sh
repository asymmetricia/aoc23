#!/bin/sh

FORCE=""
while [ "$#" -gt 2 ]; do
  case "$1" in
    -f|--force) FORCE="$1"; shift ;;
    *) break ;;
  esac
done

if [ "$#" -ne 2 ]; then
  echo "usage: init.sh YEAR DAY"
  exit 1
fi

YEAR="$(printf "%04d" "$1")"
DAY="$(printf "%02d" "$2")"
DEC_DAY="$(printf "%d" "$2")"

render() {
  SAVE=""
  if [ "$1" = "-s" ]; then
    [ -z "$FORCE" ] && SAVE="$1"
    shift
  fi

  TEMPLATE="$1"; shift
  OUTPUT="$1"; shift

  mkdir -p "$(dirname "$OUTPUT")"
  if [ -z "$SAVE" ] || ! [ -f "$OUTPUT" ]; then
    sed "s/__YEAR__/$YEAR/g; s/__DAY__/$DAY/g; s/__DEC_DAY__/$DEC_DAY/g; s/__PART__/$PART/g;" < "$TEMPLATE" > "$OUTPUT"
    git add "$OUTPUT"
  fi
}

render -s "template.go.tmpl"      "$YEAR/day$DAY/main.go"
render -s "template_test.go.tmpl" "$YEAR/day$DAY/main_test.go"
render    "run_config.xml.tmpl"   ".idea/runConfigurations/${YEAR}_${DAY}.xml"
render    "test_config.xml.tmpl"  ".idea/runConfigurations/${YEAR}_day_${DAY}_test.xml"
