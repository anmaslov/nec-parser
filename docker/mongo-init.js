db.createUser(
    {
        user: "root",
        pwd: "example",
        roles: [
            {
                role: "readWrite",
                db: "smdr"
            }
        ]
    }
);
// insert one example phone station
// write your ip address phone station
db.phones.insertOne({"ip": "127.0.0.1", "port": "8080", "enabled": true});
