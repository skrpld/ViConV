module.exports = {
  async up(db, client) {
      await db.createCollection("users", {
          validator:{
              $jsonSchema: {
                  bsonType: "object",
                  required: ["email", "password_hash", "refresh_token", "refresh_token_expiry_time"],
                  properties: {
                      email: {
                          bsonType: "string",
                      },
                      password_hash: {
                          bsonType: "string",
                      },
                      refresh_token: {
                          bsonType: "string",
                      },
                      refresh_token_expiry_time:{
                          bsonType: "date",
                      }
                  }
              }
          }
      })
  },

  async down(db, client) {
      await db.dropDatabase()
  }
};
