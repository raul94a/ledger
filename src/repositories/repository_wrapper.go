package repositories

/*
*	Registry of repositories
*
*
*/
type RepositoryWrapper struct {
	AccountRepository AccountRepository
	ClientRepository ClientRepository
	TransactionRepository TransactionRepository
}