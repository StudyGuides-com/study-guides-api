// scripts/generate_jwt.js
const jwt = require("jsonwebtoken");

const token = jwt.sign(
  {
    name: "Bruce Stockwell",
    email: "bruce.stockwell@gmail.com",
    roles: ["user", "admin"],
  },
  "kPxLHszUPxuBwSJeq9JlFdB0CQxyCT2zF0XCGITThZY=", // Replace this with NEXTAUTH_SECRET in dev
  { expiresIn: "1d" }
);

console.log(token);
