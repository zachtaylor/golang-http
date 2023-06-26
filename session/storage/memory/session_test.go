package session_memory_test

import (
	"taylz.io/http/session"
	session_memory "taylz.io/http/session/storage/memory"
)

func ServiceIsManager(s *session_memory.Service) session.Manager { return s }
