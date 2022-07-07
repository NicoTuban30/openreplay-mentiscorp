package cache

import (
	"errors"
	. "openreplay/backend/pkg/db/types"
	. "openreplay/backend/pkg/messages"
)

func (c *PGCache) InsertWebSessionStart(sessionID uint64, s *SessionStart) error {
	return c.Conn.InsertSessionStart(sessionID, &Session{
		SessionID:      sessionID,
		Platform:       "web",
		Timestamp:      s.Timestamp,
		ProjectID:      uint32(s.ProjectID),
		TrackerVersion: s.TrackerVersion,
		RevID:          s.RevID,
		UserUUID:       s.UserUUID,
		UserOS:         s.UserOS,
		UserOSVersion:  s.UserOSVersion,
		UserDevice:     s.UserDevice,
		UserCountry:    s.UserCountry,
		// web properties (TODO: unite different platform types)
		UserAgent:            s.UserAgent,
		UserBrowser:          s.UserBrowser,
		UserBrowserVersion:   s.UserBrowserVersion,
		UserDeviceType:       s.UserDeviceType,
		UserDeviceMemorySize: s.UserDeviceMemorySize,
		UserDeviceHeapSize:   s.UserDeviceHeapSize,
		UserID:               &s.UserID,
	})
}

func (c *PGCache) HandleWebSessionStart(sessionID uint64, s *SessionStart) error {
	if c.sessions[sessionID] != nil {
		return errors.New("This session already in cache!")
	}
	c.sessions[sessionID] = &Session{
		SessionID:      sessionID,
		Platform:       "web",
		Timestamp:      s.Timestamp,
		ProjectID:      uint32(s.ProjectID),
		TrackerVersion: s.TrackerVersion,
		RevID:          s.RevID,
		UserUUID:       s.UserUUID,
		UserOS:         s.UserOS,
		UserOSVersion:  s.UserOSVersion,
		UserDevice:     s.UserDevice,
		UserCountry:    s.UserCountry,
		// web properties (TODO: unite different platform types)
		UserAgent:            s.UserAgent,
		UserBrowser:          s.UserBrowser,
		UserBrowserVersion:   s.UserBrowserVersion,
		UserDeviceType:       s.UserDeviceType,
		UserDeviceMemorySize: s.UserDeviceMemorySize,
		UserDeviceHeapSize:   s.UserDeviceHeapSize,
		UserID:               &s.UserID,
	}
	if err := c.Conn.HandleSessionStart(sessionID, c.sessions[sessionID]); err != nil {
		c.sessions[sessionID] = nil
		return err
	}
	return nil
}

func (c *PGCache) InsertWebSessionEnd(sessionID uint64, e *SessionEnd) error {
	return c.InsertSessionEnd(sessionID, e.Timestamp)
}

func (c *PGCache) HandleWebSessionEnd(sessionID uint64, e *SessionEnd) error {
	return c.HandleSessionEnd(sessionID)
}

func (c *PGCache) InsertWebErrorEvent(sessionID uint64, e *ErrorEvent) error {
	session, err := c.GetSession(sessionID)
	if err != nil {
		return err
	}
	if err := c.Conn.InsertWebErrorEvent(sessionID, session.ProjectID, e); err != nil {
		return err
	}
	session.ErrorsCount += 1
	return nil
}

func (c *PGCache) InsertWebFetchEvent(sessionID uint64, e *FetchEvent) error {
	session, err := c.GetSession(sessionID)
	if err != nil {
		return err
	}
	project, err := c.GetProject(session.ProjectID)
	if err != nil {
		return err
	}
	return c.Conn.InsertWebFetchEvent(sessionID, project.SaveRequestPayloads, e)
}

func (c *PGCache) InsertWebGraphQLEvent(sessionID uint64, e *GraphQLEvent) error {
	session, err := c.GetSession(sessionID)
	if err != nil {
		return err
	}
	project, err := c.GetProject(session.ProjectID)
	if err != nil {
		return err
	}
	return c.Conn.InsertWebGraphQLEvent(sessionID, project.SaveRequestPayloads, e)
}