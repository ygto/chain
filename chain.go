package chain

type ChainParamType map[string]interface{}

type Chain struct {
	ok      bool
	next    *Chain
	preview *Chain
	params  ChainParamType
	process func(ChainParamType) (interface{}, bool)
	recover func(ChainParamType, interface{}) bool
}

func NewChain(params ChainParamType) *Chain {
	c := new(Chain)
	c.ok = false
	c.params = params

	return c
}


func (c *Chain) Chain(newChain *Chain) *Chain {
	c.next = newChain
	newChain.preview = c

	return newChain
}

func (c *Chain) SetProcess(process func(ChainParamType) (interface{}, bool)) {
	c.process = process
}

func (c *Chain) SetRecover(recover func(ChainParamType, interface{}) bool) {
	c.recover = recover
}

func (c *Chain) RunRecover(params ChainParamType, result interface{}) bool {
	if c.recover != nil {
		ok := c.recover(params, result)
		if ok && c.preview != nil {
			return c.preview.RunRecover(params, result)
		}
	}

	return false
}

func (c *Chain) Run() (interface{}, bool) {
	result, ok := c.process(c.params)

	if ok {

		//preview recover
		if c.preview != nil && c.preview.recover != nil {
			fmt.Println("RECOVER")
			c.preview.RunRecover(c.params, result)
		}

		return result, ok
	}

	if c.next != nil {

		return c.next.Run()
	}

	return nil, false
}
