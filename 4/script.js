// @ts-check

/** @typedef {[number, number]} coords */

/** @type {coords[]} */
const progression = [
  [0, -1],
  [1, -1],
  [1, 0],
  [1, 1],
  [0, 1],
  [-1, 1],
  [-1, 0],
  [-1, -1],
];

/**
 * @param {string} data
 * @return {string[][]}
 */
function stringToArray2d(data) {
  return data
    .split("\n")
    .filter(Boolean)
    .map((row) => row.split(""));
}

/**
 * @param {string[][]} array2d
 * @param {string} target
 * @param {coords[]} path
 * @param {coords} direction
 * @returns {coords[] | null}
 */
function matchInDirection(array2d, target, path, direction) {
  /** @type {coords} */
  const nextStep = [
    path[path.length - 1][0] + direction[0],
    path[path.length - 1][1] + direction[1],
  ];
  const element = array2d[nextStep[0]]?.[nextStep[1]];
  if (element === undefined) {
    return null;
  }
  if (element !== target[path.length]) {
    return null;
  }
  path.push(nextStep);
  if (target.length === path.length) {
    return path;
  }
  return matchInDirection(array2d, target, path, direction);
}

/**
 * @param {string[][]} array2d
 * @param {string} target
 * @returns {coords[][]}
 */
function seekInArray2d(array2d, target) {
  /** @type coords[][] */
  const matches = [];
  for (let x = 0; x < array2d.length; x++) {
    for (let y = 0; y < array2d.length; y++) {
      if (array2d[x][y] !== target[0]) {
        continue;
      }
      for (const delta of progression) {
        const res = matchInDirection(array2d, target, [[x, y]], delta);
        if (res) {
          matches.push(res);
        }
      }
    }
  }
  return matches;
}

/**
 * @param {string} input
 * @param {string} target
 * @returns {number}
 */
function countStringMatches(input, target) {
  const array2d = stringToArray2d(input);
  const matches = seekInArray2d(array2d, target);
  return matches.length;
}

////////PART TWO X-MAS MY ASS WTF///////////////
const compare = { M: "S", S: "M" };
const diagonal = [
  [
    [-1, -1],
    [1, 1],
  ],
  [
    [-1, 1],
    [1, -1],
  ],
];
/**
 * @param {string[][]} array2d
 * @param {coords} mid
 * @returns {boolean}
 */
function hasXmasMatch(array2d, mid) {
  let ok = true;
  for (const pair of diagonal) {
    let prev = "";
    for (let i = 0; i < pair.length; i++) {
      const [x, y] = pair[i];
      const cur = compare[array2d[mid[0] + x][mid[1] + y]];
      if ((prev && cur === prev) || cur === undefined) {
        return false;
      }
      prev = cur;
    }
  }
  return ok;
}

/**
 * @param {string} input
 * @retunrs {number}
 */
function partTwoElfWtf(input) {
  let matches = 0;
  const array2d = stringToArray2d(input);
  for (let x = 1; x < array2d.length - 1; x++) {
    for (let y = 1; y < array2d.length - 1; y++) {
      if (array2d[x][y] !== "A") {
        continue;
      }
      if (hasXmasMatch(array2d, [x, y])) {
        matches++;
      }
    }
  }
  return matches;
}
////////////////////////////////////////////////

// @ts-expect-error shutup
process.stdin.on("data", (/** @type {Buffer} */ buffer) => {
  console.log("part 1:", countStringMatches(buffer.toString(), "XMAS"));
  console.log("part 2:", partTwoElfWtf(buffer.toString()));
});
