package main

import (
	"time"
)

type RateService struct {
	MIN_TIME time.Time
	MAX_TIME time.Time
	db       *DbInterface
}

func newRateService(dbintf *DbInterface) (*RateService) {
	service := RateService{db:dbintf}
	service.MIN_TIME = time.Time{}
	service.MAX_TIME, _ = time.Parse("2006-01-02T15:04:05", "9999-12-31T23:59:59")
	return &service
}

func (s *RateService) store(rate *RateRecord) {
	log.Debug("Storing rate %s%s", rate.FromCode, rate.ToCode)
	tx := s.db.db.Begin()
	err := s.setRate(rate)
	if err != nil {
		log.Error("Unable to store %s%s rate: %v", rate.FromCode, rate.ToCode, err)
		tx.Rollback()
		return
	}
	tx.Commit()
}

func (s *RateService) setRate(rate *RateRecord) (error) {
	rate.Beginning = time.Now()
	rate.End = rate.Beginning.Add(time.Hour)

	exterior := s.db.FindExterior(rate)
	if exterior == nil {
		//First rate ever
		rate.Beginning = s.MIN_TIME
		rate.End = s.MAX_TIME
		log.Warning("Creating first rate for %s%s", rate.FromCode, rate.ToCode)
		return s.db.SaveRate(rate)
	}
	log.Info("Found exterior")
	next := s.db.FindFollowing(exterior)
	exterior.End = rate.Beginning
	err := s.db.UpdateRate(exterior)
	if err != nil {
		return err
	}

	if next == nil {
		rate.End = s.MAX_TIME
		log.Info("Set last rate for %s%s", rate.FromCode, rate.ToCode)
		return s.db.SaveRate(rate)
	}

	next.Beginning = rate.End
	err = s.db.UpdateRate(next)
	if err != nil {
		return err
	}

	log.Info("Set intermediate rate for %s%s", rate.FromCode, rate.ToCode)
	return s.db.SaveRate(rate)
	return nil
}
