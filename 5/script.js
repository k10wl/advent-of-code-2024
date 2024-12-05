// @ts-check

/** @typedef {{left: number, right: number}} PageOrderingRule */

/**
 * @param {string} input
 * @returns {{pageOrderingRules: PageOrderingRule[], production: number[][]}}
 */
function parseInput(input) {
  const [rulesString, productionString] = input.split("\n\n").filter(Boolean);
  const pageOrderingRules = rulesString.split("\n").map((row) => {
    const [left, right] = row.split("|").filter(Boolean);
    return { left: +left, right: +right };
  });
  const production = productionString
    .split("\n")
    .filter(Boolean)
    .map((row) => {
      return row
        .split(",")
        .filter(Boolean)
        .map((string) => +string);
    });
  return { pageOrderingRules, production };
}

/** @type {Map<number, number[]>} */
const rules = new Map();
/** @param {PageOrderingRule} pageOrderingRules */
function mapRule(pageOrderingRules) {
  let stored = rules.get(pageOrderingRules.right);
  if (!stored) {
    stored = [];
    rules.set(pageOrderingRules.right, stored);
  }
  stored.push(pageOrderingRules.left);
}

/**
 * @param {number[]} production
 * @param {Map<number, number[]>} rules
 * @param {boolean} fixed
 * @return {{ordered: number[], fixed: boolean}}
 */
function checkOrder(production, rules, fixed) {
  const result = { ordered: production, fixed };
  const indexesMap = result.ordered.reduce((map, value, index) => {
    map[value] = index;
    return map;
  }, /** @type {Record<number, number>} */ ({}));
  for (let i = 0; i < result.ordered.length; i++) {
    const page = result.ordered[i];
    const requirements = rules.get(page);
    if (!requirements) {
      continue;
    }
    for (const requirement of requirements) {
      const storedRequirement = indexesMap[requirement];
      if (storedRequirement === undefined) {
        continue;
      }
      if (storedRequirement > i) {
        result.ordered[i] = result.ordered[storedRequirement];
        result.ordered[storedRequirement] = page;
        return checkOrder(result.ordered, rules, true);
      }
    }
  }
  return result;
}

// @ts-ignore shutup
process.stdin.on("data", (/** @type {Buffer} */ data) => {
  const input = parseInput(data.toString());
  input.pageOrderingRules.forEach(mapRule);
  let sum = 0;
  input.production.forEach((line) => {
    const result = checkOrder([...line], rules, false);
    if (result.fixed) {
      sum += result.ordered[Math.floor(result.ordered.length / 2)];
      return;
    } else {
      // sum += result.ordered[Math.floor(result.ordered.length / 2)];
    }
  });
  console.log(sum);
});
