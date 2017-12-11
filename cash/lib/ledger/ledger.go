package ledger

import (
  "time"
  "fmt"

  "strings"

)

const (
  DateFmt = "2006/01/02"
)

type AccountName []string

func (a AccountName) String() string {
  return strings.Join(a, ":")
}


func Liability(pieces ...string) AccountName{
  return append([]string{"liabilities"}, pieces...)
}

func Asset(pieces ...string) AccountName{
  return append([]string{"assets"}, pieces...)
}

func Expense(pieces ...string) AccountName{
  return append([]string{"expenses"}, pieces...)
}

type Change struct {
  Account AccountName
  Amount float64
}

func (c *Change) String() string {
  if c.Amount == 0 {
    return c.Account.String()
  }
  return fmt.Sprintf("%s    %.2f", c.Account.String(), c.Amount)
}

type Transaction struct {
  Date time.Time
  Description string
  Changes []Change
}

func (t *Transaction) String() string {
  lines := make([]string, 1 + len(t.Changes))

  lines[0] = fmt.Sprintf("%s %s", t.Date.Format(DateFmt),
                           t.Description)

  for i, c := range t.Changes {
    lines[i+1] = "    " + c.String()
  }

  return strings.Join(lines, "\n")
}
