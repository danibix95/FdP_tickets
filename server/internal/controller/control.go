package controller

import (
	"fmt"
	"github.com/danibix95/FdP_tickets/server/internal/dbconn"
	"github.com/kataras/iris"
	"io"
	"log"
	"os"
	"time"
)

type AppController struct {
	dbc    *dbconn.DBController
	logger *log.Logger
}

func New(controlLogFile *os.File, dbLogFile *os.File) *AppController {
	appc := AppController{
		dbc: dbconn.New(dbLogFile),
		logger: log.New(io.MultiWriter(os.Stderr, controlLogFile),
			"", log.LstdFlags),
	}

	return &appc
}

func (appc *AppController) Ping(ctx iris.Context) {
	appc.dbc.PingDB() // test db connection

	_, _ = ctx.JSON(iris.Map{
		"message": fmt.Sprintf("Pong - %v", time.Now().Local()),
	})
}

/* ======= LOGIN MANAGEMENT ======= */
func (appc *AppController) RequireLogin(ctx iris.Context) {
	appc.logger.Println("User interaction with app content!")
	if false {
		ctx.StatusCode(iris.StatusUnauthorized)
		ctx.EndRequest()
	} else {
		// move forward the execution of request chain
		ctx.Next()
	}
}

func (appc *AppController) Login(ctx iris.Context) {
}

func (appc *AppController) Logout(ctx iris.Context) {
}

func (appc *AppController) IsAdmin(ctx iris.Context) {
}

/* ======= TICKETS MANAGEMENT ======= */
// Used to check whether and when an attendee entered to the party
func (appc *AppController) WhenEntered(ctx iris.Context) {
	ticketNum, err := ctx.Params().GetUint("ticketNum")
	if err != nil {
		appc.logger.Panicln(fmt.Sprintf("Ticket number %v is not valid!"+
			" Ticket numbers should be a natural number.", ticketNum))
	}
	enteredTime := appc.dbc.WhenEntered(ticketNum)

	_, _ = ctx.JSON(iris.Map{
		"ticketNum": ticketNum,
		"time":      enteredTime.Time,
		"isEntered": enteredTime.Valid,
		"status":    200,
	})
}

func (appc *AppController) GetTickets(ctx iris.Context) {
}

func (appc *AppController) GetTicketsInfo(ctx iris.Context) {
}

func (appc *AppController) GetTicketDetails(ctx iris.Context) {
}

func (appc *AppController) SetEntered(ctx iris.Context) {
	var ticket struct {
		TicketNum uint `json:"ticketNum"`
	}
	err := ctx.ReadJSON(&ticket)
	if err != nil {
		appc.logger.Panicln(fmt.Sprintf("Error reading JSON!\n%v", err))
	}

	// check ticket is in range
	if ticket.TicketNum > dbconn.TICKETHIGH {
		appc.logger.Panicln(fmt.Sprintf("Ticket number %v is not valid!"+
			" Ticket number is not in specified range.", ticket.TicketNum))
	}

	// Default result -> the ticket has not been sold
	// and therefore it cannot enter without first pay the entrance
	result := iris.Map{
		"ticketNum": ticket.TicketNum,
		"status":    400,
		"entered":   false,
		"msg":       "Ticket unsold!",
	}

	switch appc.dbc.IsSoldEntered(ticket.TicketNum) {
	case dbconn.SOLDENTERED:
		// notify that the ticket is already entered,
		// so the same ticket number cannot enter again
		result["entered"] = true
		result["msg"] = "Ticket sold and already entered."
	case dbconn.SOLD:
		if appc.dbc.SetEntered(ticket.TicketNum) {
			// successful update
			result["status"] = 200
			result["entered"] = true
			result["msg"] = "Ticket entered correctly!"
		} else {
			result["status"] = 500
			result["msg"] = "Error encountered while allowing the entrance to this ticket..."
		}
	}

	_, _ = ctx.JSON(result)
}

func (appc *AppController) ConfirmEntrance(ctx iris.Context) {
}

func (appc *AppController) SellTicket(ctx iris.Context) {
}

func (appc *AppController) RollbackEntrance(ctx iris.Context) {
}

func (appc *AppController) GetTicketVendor(ctx iris.Context) {
}

/* ======= ERROR MANAGEMENT ======= */
func (appc *AppController) Unauthorized(ctx iris.Context) {
	_, _ = ctx.JSON(iris.Map{
		"status": 401,
	})
}

func (appc *AppController) NotFound(ctx iris.Context) {
	_, _ = ctx.JSON(iris.Map{
		"status": 404,
	})
}
