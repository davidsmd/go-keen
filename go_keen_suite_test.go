package keen_test

import (
	"log"

	"github.com/joho/godotenv"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestGoKeen(t *testing.T) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	RegisterFailHandler(Fail)
	RunSpecs(t, "GoKeen Suite")
}
