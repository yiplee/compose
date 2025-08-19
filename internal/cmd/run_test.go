package cmd

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/spf13/viper"
)

func BenchmarkLoadTasks(b *testing.B) {
	// Setup test configuration
	viper.SetConfigType("yaml")
	testConfig := `
tasks:
  service_1:
    cmds: sleep 5
  service_2:
    cmds: sleep 5
    delay: 1s
  service_3:
    cmds: echo "hello"
  service_4:
    cmds: ls -la
  service_5:
    cmds: pwd
`
	
	// Create a temporary viper instance for testing
	testViper := viper.New()
	testViper.SetConfigType("yaml")
	testViper.ReadConfig(bytes.NewReader([]byte(testConfig)))
	
	// Store original viper instance
	originalViper := viper.GetViper()
	defer func() {
		// Restore original viper
		viper.SetConfigFile(originalViper.ConfigFileUsed())
	}()
	
	// Mock the viper.Sub call
	viper.Set("tasks", testViper.Get("tasks"))
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		loadTasks()
	}
}

func BenchmarkLoadTasksLarge(b *testing.B) {
	// Setup large test configuration
	viper.SetConfigType("yaml")
	
	// Generate 100 test tasks
	testConfig := "tasks:\n"
	for i := 0; i < 100; i++ {
		testConfig += fmt.Sprintf("  service_%d:\n    cmds: sleep %d\n", i, i%10+1)
	}
	
	// Create a temporary viper instance for testing
	testViper := viper.New()
	testViper.SetConfigType("yaml")
	testViper.ReadConfig(bytes.NewReader([]byte(testConfig)))
	
	// Store original viper instance
	originalViper := viper.GetViper()
	defer func() {
		// Restore original viper
		viper.SetConfigFile(originalViper.ConfigFileUsed())
	}()
	
	// Mock the viper.Sub call
	viper.Set("tasks", testViper.Get("tasks"))
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		loadTasks()
	}
}