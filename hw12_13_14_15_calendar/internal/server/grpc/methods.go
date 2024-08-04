package grpc

import (
	"context"
	"time"

	pb "github.com/ovs325/ovs-otus/hw12_13_14_15_calendar/grpc/event_service"
	cm "github.com/ovs325/ovs-otus/hw12_13_14_15_calendar/internal/common"
	er "github.com/ovs325/ovs-otus/hw12_13_14_15_calendar/internal/errors"
	tp "github.com/ovs325/ovs-otus/hw12_13_14_15_calendar/internal/types"
)

func (s *GrpcServer) CreateEvent(ctx context.Context, req *pb.BodyEventRequest) (response *pb.CreateEventResponse, err error) {
	response_ := pb.CreateEventResponse{Id: int64(0), Success: false}
	checkItem, err := s.bodyToCheckItem(req)
	if err != nil {
		s.log.Error("ошибка grpc-клиента", "error", err.Error())
		return &response_, err
	}
	id, err := s.logic.CreateEventLogic(ctx, &checkItem)
	if err != nil {
		s.log.Error("ошибка grpc-сервера", "error", err.Error())
		return &response_, err
	}
	response_.Id = int64(id)
	response_.Success = true
	return &response_, err
}

func (s *GrpcServer) UpdateEvent(ctx context.Context, req *pb.BodyEventRequest) (response *pb.UpdateEventResponse, err error) {
	response_ := pb.UpdateEventResponse{Success: false}
	checkItem, err := s.bodyToCheckItem(req)
	if err != nil {
		s.log.Error("ошибка grpc-клиента", "error", err.Error())
		return &response_, err
	}
	if err = s.logic.UpdateEventLogic(ctx, &checkItem); err != nil {
		s.log.Error("ошибка grpc-сервера", "error", err.Error())
		return &response_, err
	}
	response_.Success = true
	return &response_, err
}

func (s *GrpcServer) DeleteEvent(ctx context.Context, req *pb.DeleteEventRequest) (response *pb.DeleteEventResponse, err error) {
	response_ := pb.DeleteEventResponse{Success: false}
	if err = s.logic.DelEventLogic(ctx, req.Id); err != nil {
		s.log.Error("ошибка grpc-сервера", "error", err.Error())
		return &response_, err
	}
	response_.Success = true
	return
}

func (s *GrpcServer) GetEventsDay(ctx context.Context, req *pb.GetEventsRequest) (response *pb.GetEventsResponse, err error) {
	response_ := pb.GetEventsResponse{Success: false}
	date, err := time.Parse(time.RFC3339, req.Date)
	if err != nil {
		s.log.Error("ошибка grpc-клиента", "error", err.Error())
		return &response_, er.ErrBadFormatTime
	}
	res, err := s.logic.GetDayLogic(ctx, date, cm.Paginate{})
	if err != nil {
		s.log.Error("ошибка grpc-сервера", "error", err.Error())
		return &response_, err
	}
	response_.Events = s.eventsToResponce(res)
	response_.Success = true
	return &response_, nil
}

func (s *GrpcServer) GetEventsMonth(ctx context.Context, req *pb.GetEventsRequest) (response *pb.GetEventsResponse, err error) {
	response_ := pb.GetEventsResponse{Success: false}
	date, err := time.Parse(time.RFC3339, req.Date)
	if err != nil {
		s.log.Error("ошибка grpc-клиента", "error", err.Error())
		return &response_, er.ErrBadFormatTime
	}
	res, err := s.logic.GetWeekLogic(ctx, date, cm.Paginate{})
	if err != nil {
		s.log.Error("ошибка grpc-сервера", "error", err.Error())
		return &response_, err
	}
	response_.Events = s.eventsToResponce(res)
	response_.Success = true
	return &response_, nil
}

func (s *GrpcServer) GetEventsWeek(ctx context.Context, req *pb.GetEventsRequest) (response *pb.GetEventsResponse, err error) {
	response_ := pb.GetEventsResponse{Success: false}
	date, err := time.Parse(time.RFC3339, req.Date)
	if err != nil {
		s.log.Error("ошибка grpc-клиента", "error", err.Error())
		return &response_, er.ErrBadFormatTime
	}
	res, err := s.logic.GetMonthLogic(ctx, date, cm.Paginate{})
	if err != nil {
		s.log.Error("ошибка grpc-сервера", "error", err.Error())
		return &response_, err
	}
	response_.Events = s.eventsToResponce(res)
	response_.Success = true
	return &response_, nil
}

func (s *GrpcServer) bodyToCheckItem(req *pb.BodyEventRequest) (checkItem tp.EventRequest, err error) {
	checkItem = tp.EventRequest{
		Event: tp.Event{ID: req.Event.Id,
			Name:        req.Event.Name,
			Description: req.Event.Description,
			UserID:      req.Event.UserId},
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

func (s *GrpcServer) eventsToResponce(res tp.QueryPage[tp.EventModel]) (eventModelList []*pb.EventModel) {
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
