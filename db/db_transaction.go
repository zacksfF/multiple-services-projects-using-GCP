package db

//Beginner represents a value that can begin that can begin a transaction.
type Beginner interface{
	Begin() (CommitrollBacker, error)
}

// CommitRollbacker represents a value that can commit or rollback a transaction.
type CommitrollBacker interface{
	Commit() error
	Rollback() error
}