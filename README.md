# memdb - An in-memory database

This is a simple in-memory database, which supports basic
read/write functionality, as well as simple transactions.

The easiest way to run this from source is to simply do:

```zsh
make run
```

This will open the following prompt:

```zsh
Welcome to MemDB!
>>
```

At the prompt, you can enter the following commands:

`SET [name] [value]`

  Sets the name in the database to the given value

`GET [name]`

  Prints the value for the given name. If the value is not
  in the database, prints `NULL`

`DELETE [name]`

  Deletes the value from the database

`COUNT [value]`

  Returns the number of names that have the given value
  assigned to them. If that value is not assigned
  anywhere, prints `0`

`END`

  Exits the database

`BEGIN`

  Begins a new transaction

`ROLLBACK`

  Rolls back the most recent transaction. If there is no
  transaction to rollback, prints `TRANSACTION NOT FOUND`
  error

`COMMIT`

  Commits *all* of the open transactions

If you get stuck, just type `HELP` to get the list of
commands displayed to you.

The actual commands are not case-sensitive, so you can
type `set`, `Set`, or `SET`, for example, and get the
same results.
