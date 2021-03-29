package main

type msgFlags []string

func (i *msgFlags) String() string {
	return "my string representation"
}

func (i *msgFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}
