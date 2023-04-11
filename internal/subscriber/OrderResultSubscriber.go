package subscriber

type OrderResultSubscriber interface {
	HandleOrderResult() error
}
