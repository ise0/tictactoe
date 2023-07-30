const d = new Date();
const h = d.getHours().toString().padStart(2, "0");
const m = d.getMinutes().toString().padStart(2, "0");
const s = d.getSeconds().toString().padStart(2, "0");
const ms = d.getMilliseconds().toString().padStart(1, "0");
console.log(ms);
console.log(`${h}:${m}:${s},${ms}`);