/**
 * @param {number[]} a
 * @param {number[]} b
 * @return {number}
 */
function findDistance(a, b) {
  const sortedA = sortArr(a);
  const sortedB = sortArr(b);
  let distance = 0;
  for (let i = 0; i < sortedA.length; i++) {
    const numA = sortedA[i];
    const numB = sortedB[i];
    if (numA === undefined || numB === undefined) {
      throw new Error("bad data");
    }
    distance += Math.abs(numA - numB);
  }
  return distance;
}

/**
 * @param {number[]} arr
 * @returns {number[]}
 */
function sortArr(arr) {
  return [...arr].sort((a, b) => a - b);
}

/** @param {string} str
 * @returns {[number[], number[]]}
 */
function parseLists(str) {
  const pairs = str.split("\n");
  const a = [];
  const b = [];
  for (const pair of pairs) {
    const [left, right] = pair.split(/\s+/);
    if (left === undefined || right === undefined) {
      continue;
    }
    a.push(+left);
    b.push(+right);
  }
  return [a, b];
}

/**
 * @param {number[]} a
 * @param {number[]} b
 * @return {number}
 */
function calcSimularityScore(a, b) {
  /** @type {Map<number, {a: number, b: number}>} */
  const simularities = new Map();
  for (const num of a) {
    simularities.set(num, { a: (simularities.get(num)?.a ?? 0) + 1, b: 0 });
  }
  for (const num of b) {
    const stored = simularities.get(num);
    simularities.set(num, {
      a: stored?.a || 0,
      b: (simularities.get(num)?.b ?? 0) + 1,
    });
  }
  return [...simularities.entries()].reduce((prev, [num, { a, b }]) => {
    return prev + num * a * b;
  }, 0);
}

process.stdin.on("data", (data) => {
  console.log("part 1:", findDistance(...parseLists(data.toString())));
  console.log("part 2:", calcSimularityScore(...parseLists(data.toString())));
});
