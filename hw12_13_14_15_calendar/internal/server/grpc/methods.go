package grpc

import (
	"context"
	"time"

	pb "github.com/ovs325/ovs-otus/hw12_13_14_15_calendar/grpc/event_service"
	cm "github.com/ovs325/ovs-otus/hw12_13_14_15_calendar/internal/common"
	er "github.com/ovs325/ovs-otus/hw12_13_14_15_calendar/internal/errors"
	tp "github.com/ovs325/ovs-otus/hw12_13_14_15_calendar/internal/types"
	"github.com/ovs325/ovs-otus/hw12_13_14_15_calendar/pkg"
)

func (s *ProtoServer) CreateEvent(ctx context.Context, req *pb.BodyEventRequest) (*pb.CreateEventResponse, error) {
	res := pb.CreateEventResponse{Id: int64(0), Success: false}
	checkItem, err := s.bodyToCheckItem(req)
	if err != nil {
		s.log.Error("ошибка grpc-клиента", "error", err.Error())
		return &res, err
	}
	checkItem.ID = 0

	event := tp.EventModel{}
	event.GetModel(checkItem)

	id, err := s.stor.CreateEvent(ctx, &event)
	if err != nil {
		s.log.Error("ошибка grpc-сервера", "error", err.Error())
		return &res, err
	}
	res.Id = id
	res.Success = true
	return &res, err
}

func (s *ProtoServer) UpdateEvent(ctx context.Context, req *pb.BodyEventRequest) (*pb.UpdateEventResponse, error) {
	res := pb.UpdateEventResponse{Success: false}
	checkItem, err := s.bodyToCheckItem(req)
	if err != nil {
		s.log.Error("ошибка grpc-клиента", "error", err.Error())
		return &res, err
	}
	if checkItem.ID == 0 {
		s.log.Error("ошибка grpc-клиента: id не должен быть <= 0")
		return &res, err
	}
	event := tp.EventModel{}
	event.GetModel(checkItem)
	if err = s.stor.UpdateEvent(ctx, &event); err != nil {
		s.log.Error("ошибка grpc-сервера", "error", err.Error())
		return &res, err
	}
	res.Success = true
	return &res, err
}

func (s *ProtoServer) DeleteEvent(ctx context.Context, req *pb.DeleteEventRequest) (*pb.DeleteEventResponse, error) {
	res := pb.DeleteEventResponse{Success: false}
	if err := s.stor.DelEvent(ctx, req.Id); err != nil {
		s.log.Error("ошибка grpc-сервера", "error", err.Error())
		return &res, err
	}
	res.Success = true
	return &res, nil
}

func (s *ProtoServer) GetEventsDay(ctx context.Context, req *pb.GetEventsRequest) (*pb.GetEventsResponse, error) {
	return s.getEventsInterval(ctx, "day", req)
}

func (s *ProtoServer) GetEventsWeek(ctx context.Context, req *pb.GetEventsRequest) (*pb.GetEventsResponse, error) {
	return s.getEventsInterval(ctx, "week", req)
}

func (s *ProtoServer) GetEventsMonth(ctx context.Context, req *pb.GetEventsRequest) (*pb.GetEventsResponse, error) {
	return s.getEventsInterval(ctx, "month", req)
}

func (s *ProtoServer) getEventsInterval(
	ctx context.Context,
	interval string,
	req *pb.GetEventsRequest,
) (*pb.GetEventsResponse, error) {
	res := pb.GetEventsResponse{Success: false}
	date, err := time.Parse(time.RFC3339, req.Date)
	if err != nil {
		s.log.Error("ошибка grpc-клиента", "error", err.Error())
		return &res, er.ErrBadFormatTime
	}
	var first, last time.Time
	switch interval {
	case "day":
		first, last = pkg.GetDayInterval(date)
	case "week":
		first, last = pkg.GetWeekInterval(date)
	case "month":
		first, last = pkg.GetMonthInterval(date)
	}
	result, err := s.stor.GetEventsForTimeInterval(ctx, first, last, cm.Paginate{})
	if err != nil {
		s.log.Error("ошибка grpc-сервера", "error", err.Error())
		return &res, err
	}
	res.Events = s.eventsToResponce(result)
	res.Success = true
	return &res, nil
}

func (s *ProtoServer) bodyToCheckItem(req *pb.BodyEventRequest) (checkItem tp.EventRequest, err error) {
	checkItem = tp.EventRequest{
		Event: tp.Event{
			ID:          req.Event.Id,
			Name:        req.Event.Name,
			Description: req.Event.Description,
			UserID:      req.Event.UserId,
		},
		NDayAlarm: int(req.NDayAlarm),
	}
	if req.Event.Date == "" {
		checkItem.Date = time.Now()
	} else {
		checkItem.Date, err = time.Parse(time.RFC3339, req.Event.Date)
		if err != nil {
			return tp.EventRequest{}, er.ErrBadFormatTime
		}
	}
	if req.Event.Expiry == "" {
		checkItem.Expiry = checkItem.Date.AddDate(0, 0, 1)
	} else {
		checkItem.Expiry, err = time.Parse(time.RFC3339, req.Event.Expiry)
		if err != nil {
			return tp.EventRequest{}, er.ErrBadFormatTime
		}
	}
	return
}

func (s *ProtoServer) eventsToResponce(res tp.QueryPage[tp.EventModel]) (eventModelList []*pb.EventModel) {
	format := "2006-01-02 15:04:05.00"
	eventModelList = make([]*pb.EventModel, 0, len(res.Content))
	for _, event := range res.Content {
		model := pb.EventModel{
			Event: &pb.Event{
				Id:          event.ID,
				Name:        event.Name,
				Date:        event.Date.Format(format),
				Expiry:      event.Expiry.Format(format),
				Description: event.Description,
				UserId:      event.UserID,
			},
			TimeAlarm: event.TimeAlarm.Format(format),
		}
		eventModelList = append(eventModelList, &model)
	}
	return eventModelList
}
