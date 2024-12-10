const actions = [
  (numA, numB) => numA + numB,
  (numA, numB) => numA * numB,
  (numA, numB) => +`${numA}${numB}`,
];

/**
 * @param {number} initial
 * @param {number[]} data
 * @param {number} limit
 * @returns {number}
 */
function traverce(initial, data, limit) {
  if (data.length === 0) {
    return initial;
  }

  for (const action of actions) {
    const current = action(initial, data[0]);
    if (limit < current) {
      continue;
    }
    let res = traverce(current, data.slice(1), limit);
    if (res === limit) {
      return res;
    }
  }

  return -1;
}

// @ts-ignore shut up
process.stdin.on("data", (/** @type {Buffer} */ data) => {
  console.time();
  const input = parseInput(data.toString());
  let sum = 0;
  for (const { target, values } of input) {
    const value = traverce(values[0], values.slice(1), target);
    if (value === target) {
      sum += value;
    }
  }
  console.log(sum);
  console.timeEnd();
});

/**
 * @param {string} input
 */
function parseInput(input) {
  return input
    .split("\n")
    .filter(Boolean)
    .map((row) => {
      const [target, strings] = row.split(": ").filter(Boolean);
      return {
        target: +target,
        values: strings
          .split(" ")
          .filter(Boolean)
          .map((str) => +str),
      };
    });
}
