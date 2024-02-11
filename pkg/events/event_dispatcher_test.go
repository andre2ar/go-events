package events

import (
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"sync"
	"testing"
	"time"
)

type TestEvent struct {
	Name    string
	Payload any
}

func (e *TestEvent) GetDateTime() time.Time {
	//TODO implement me
	panic("implement me")
}

func (e *TestEvent) GetName() string {
	return e.Name
}

func (e *TestEvent) GetPayload() any {
	return e.Payload
}

type TestEventHandler struct {
	ID int
}

func (h *TestEventHandler) Handle(event EventInterface, wg *sync.WaitGroup) {
	wg.Done()
}

type EventDispatcherTestSuite struct {
	suite.Suite
	event1          TestEvent
	event2          TestEvent
	handler1        TestEventHandler
	handler2        TestEventHandler
	handler3        TestEventHandler
	eventDispatcher *EventDispatcher
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(EventDispatcherTestSuite))
}

func (suite *EventDispatcherTestSuite) SetupTest() {
	suite.eventDispatcher = NewEventDispatcher()
	suite.handler1 = TestEventHandler{
		ID: 1,
	}
	suite.handler2 = TestEventHandler{
		ID: 2,
	}
	suite.handler3 = TestEventHandler{
		ID: 3,
	}
	suite.event1 = TestEvent{
		Name:    "TestEvent1",
		Payload: "test1",
	}
	suite.event2 = TestEvent{
		Name:    "TestEvent2",
		Payload: "test2",
	}
}

func (suite *EventDispatcherTestSuite) TestEventDispatcher_Register() {
	err := suite.eventDispatcher.Register(suite.event1.GetName(), &suite.handler1)
	suite.NoError(err)
	suite.Equal(1, len(suite.eventDispatcher.handlers[suite.event1.GetName()]))

	err = suite.eventDispatcher.Register(suite.event1.GetName(), &suite.handler2)
	suite.NoError(err)
	suite.Equal(2, len(suite.eventDispatcher.handlers[suite.event1.GetName()]))

	suite.Equal(&suite.handler1, suite.eventDispatcher.handlers[suite.event1.GetName()][0])
	suite.Equal(&suite.handler2, suite.eventDispatcher.handlers[suite.event1.GetName()][1])
}

func (suite *EventDispatcherTestSuite) TestEventDispatcher_Register_WithSameHandler() {
	err := suite.eventDispatcher.Register(suite.event1.GetName(), &suite.handler1)
	suite.NoError(err)
	suite.Equal(1, len(suite.eventDispatcher.handlers[suite.event1.GetName()]))

	err = suite.eventDispatcher.Register(suite.event1.GetName(), &suite.handler1)
	suite.Error(err)
	suite.Equal(1, len(suite.eventDispatcher.handlers[suite.event1.GetName()]))
}

func (suite *EventDispatcherTestSuite) TestEventDispatcher_Clear() {
	_ = suite.eventDispatcher.Register(suite.event1.GetName(), &suite.handler1)
	_ = suite.eventDispatcher.Register(suite.event1.GetName(), &suite.handler2)

	_ = suite.eventDispatcher.Register(suite.event2.GetName(), &suite.handler3)

	suite.eventDispatcher.Clear()
	suite.Equal(0, len(suite.eventDispatcher.handlers))
}

func (suite *EventDispatcherTestSuite) TestEventDispatcher_Has() {
	_ = suite.eventDispatcher.Register(suite.event1.GetName(), &suite.handler1)
	_ = suite.eventDispatcher.Register(suite.event1.GetName(), &suite.handler2)

	suite.True(suite.eventDispatcher.Has(suite.event1.GetName(), &suite.handler1))
	suite.True(suite.eventDispatcher.Has(suite.event1.GetName(), &suite.handler2))
	suite.False(suite.eventDispatcher.Has(suite.event1.GetName(), &suite.handler3))

	suite.False(suite.eventDispatcher.Has(suite.event2.GetName(), &suite.handler1))
}

type MockEventHandler struct {
	mock.Mock
}

func (m *MockEventHandler) Handle(event EventInterface, wg *sync.WaitGroup) {
	wg.Done()
	m.Called(event)
}
func (suite *EventDispatcherTestSuite) TestEventDispatcher_Dispatch() {
	eventHandler := &MockEventHandler{}
	eventHandler.On("Handle", &suite.event1)

	suite.eventDispatcher.Register(suite.event1.GetName(), eventHandler)
	suite.eventDispatcher.Dispatch(&suite.event1)

	eventHandler.AssertExpectations(suite.T())
	eventHandler.AssertNumberOfCalls(suite.T(), "Handle", 1)
}

func (suite *EventDispatcherTestSuite) TestEventDispatcher_Remove() {
	_ = suite.eventDispatcher.Register(suite.event1.GetName(), &suite.handler1)
	_ = suite.eventDispatcher.Register(suite.event1.GetName(), &suite.handler2)

	_ = suite.eventDispatcher.Register(suite.event2.GetName(), &suite.handler3)

	suite.eventDispatcher.Remove(suite.event1.GetName(), &suite.handler1)
	suite.Equal(1, len(suite.eventDispatcher.handlers[suite.event1.GetName()]))
	suite.Equal(&suite.handler2, suite.eventDispatcher.handlers[suite.event1.GetName()][0])

	suite.eventDispatcher.Remove(suite.event1.GetName(), &suite.handler2)
	suite.Equal(0, len(suite.eventDispatcher.handlers[suite.event1.GetName()]))

	suite.eventDispatcher.Remove(suite.event2.GetName(), &suite.handler3)
	suite.Equal(0, len(suite.eventDispatcher.handlers[suite.event2.GetName()]))

}
