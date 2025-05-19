package transaction_test

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	cliententity "src/domain/client"
	"src/test/utils"
	"testing"
	"time"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)