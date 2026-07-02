const { Database } = require("bun:sqlite");
const db = new Database("C:/Users/zains/AppData/Local/kiro-cli/data.sqlite3");

console.log("\nSample rows from conversations_v2:");
try {
  const rows = db.query("SELECT * FROM conversations_v2 LIMIT 5").all();
  for (const r of rows) {
    console.log(JSON.stringify(r, null, 2));
  }
} catch (e) {
  console.error("Error reading conversations_v2:", e);
}

console.log("\nSample rows from history:");
try {
  const rows = db.query("SELECT * FROM history LIMIT 5").all();
  for (const r of rows) {
    console.log(JSON.stringify(r, null, 2));
  }
} catch (e) {
  console.error("Error reading history:", e);
}

db.close();
