package cli_test

import (
	"testing"

	"github.com/nholuongut/terragrunt/pkg/cli"
	"github.com/stretchr/testify/assert"
)

var mockArgs = func() cli.Args { return cli.Args{"one", "-foo", "two", "--bar", "value"} }

func TestArgsSlice(t *testing.T) {
	t.Parallel()

	actual := mockArgs().Slice()
	expected := []string(mockArgs())
	assert.Equal(t, expected, actual)
}

func TestArgsTail(t *testing.T) {
	t.Parallel()

	actual := mockArgs().Tail()
	expected := mockArgs()[1:]
	assert.Equal(t, expected, actual)
}

func TestArgsFirst(t *testing.T) {
	t.Parallel()

	actual := mockArgs().First()
	expected := mockArgs()[0]
	assert.Equal(t, expected, actual)
}

func TestArgsGet(t *testing.T) {
	t.Parallel()

	actual := mockArgs().Get(2)
	expected := "two"
	assert.Equal(t, expected, actual)
}

func TestArgsLen(t *testing.T) {
	t.Parallel()

	actual := mockArgs().Len()
	expected := 5
	assert.Equal(t, expected, actual)
}

func TestArgsPresent(t *testing.T) {
	t.Parallel()

	actual := mockArgs().Present()
	expected := true
	assert.Equal(t, expected, actual)

	args := cli.Args([]string{})
	actual = args.Present()
	expected = false
	assert.Equal(t, expected, actual)
}

func TestArgsCommandName(t *testing.T) {
	t.Parallel()

	actual := mockArgs().CommandName()
	expected := "one"
	assert.Equal(t, expected, actual)

	args := mockArgs()[1:]
	actual = args.CommandName()
	expected = ""
	assert.Equal(t, expected, actual)
}

func TestArgsNormalize(t *testing.T) {
	t.Parallel()

	actual := mockArgs().Normalize(cli.SingleDashFlag).Slice()
	expected := []string{"one", "-foo", "two", "-bar", "value"}
	assert.Equal(t, expected, actual)

	actual = mockArgs().Normalize(cli.DoubleDashFlag).Slice()
	expected = []string{"one", "--foo", "two", "--bar", "value"}
	assert.Equal(t, expected, actual)
}
