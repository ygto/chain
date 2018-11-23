package chain


type ChainParamType map[string]interface{}

type Chain struct {
	ok      bool
	next    *Chain
	preview *Chain
	process func(ChainParamType) (interface{}, bool)
	recover func(ChainParamType, interface{}) bool
}

func NewChain() *Chain {
	c := new(Chain)
	c.ok = false

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

func (c *Chain) Run(params ChainParamType) (interface{}, bool) {
	result, ok := c.process(params)

	if ok {

		//preview recover
		if c.preview != nil {
			c.preview.RunRecover(params, result)
		}

		return result, ok
	}

	if c.next != nil {

		return c.next.Run(params)
	}

	return nil, false
}
