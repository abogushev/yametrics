client_run:
	go run cmd/agent/main.go
client_trace_profile:
	go tool pprof -http=":9090" -seconds=120 http://localhost:8100/debug/pprof/profile
client_trace_profile_and_save:
	curl -sK -v 'http://localhost:8100/debug/pprof/profile?seconds=300' > profiles/client_profile.out && go tool pprof -http=":9090" profiles/client_profile.out
	
client_trace_heap:
	go tool pprof -http=":9090" -seconds=120 http://localhost:8100/debug/pprof/heap
client_trace_heap_and_save:
	curl -sK -v 'http://localhost:8100/debug/pprof/heap?seconds=300' > profiles/client_heap.out && go tool pprof -http=":9090" profiles/client_heap.out


server_run:
	go run cmd/server/main.go
server_trace_profile:
	go tool pprof -http=":9090" -seconds=120 http://localhost:8080/debug/pprof/profile
server_trace_profile_and_save:
	curl -sK -v 'http://localhost:8080/debug/pprof/profile?seconds=300' > profiles/server_profile.out && go tool pprof -http=":9090" profiles/server_profile.out

server_trace_heap:
        go tool pprof -http=":9090" -seconds=120 http://localhost:8080/debug/pprof/heap
server_trace_heap_and_save:
        curl -sK -v 'http://localhost:8080/debug/pprof/heap?seconds=300' > profiles/server_heap.out && go tool pprof -http=":9090" profiles/server_heap.out
