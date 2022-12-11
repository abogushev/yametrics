client_agent_port = 8100
server_agent_port = 8200
server_profile_out_file = server_profile.out
server_heap_out_file = server_heap.out

trace_fn = go tool pprof -http=\":9090\" -seconds=120 http://localhost:$(1)/debug/pprof/$(2)
trace_and_save_fn = curl -sK -v http://localhost:$(1)/debug/pprof/$(2)?seconds=300 > profiles/$(3) && go tool pprof -http=":9090" profiles/$(3)

proto_gen:
	protoc --go_out=. --go_opt=paths=source_relative \
      --go-grpc_out=. --go-grpc_opt=paths=source_relative \
      internal/protocol/proto/metrics.proto

client_run:
	# go run cmd/agent/main.go
	go run -ldflags "-X main.buildVersion=$$(cat cmd/agent/version.txt) -X 'main.buildDate=$$(date +'%d/%m/%Y')' -X 'main.buildCommit=$$(git rev-parse HEAD)'" cmd/agent/main.go
client_trace_profile:
	$(call trace_fn,$(client_agent_port),profile)
client_trace_profile_and_save:
	$(call trace_and_save_fn,$(client_agent_port),profile,client_profile.out)
client_trace_heap:
	$(call trace_fn,$(client_agent_port),heap)
client_trace_heap_and_save:
	$(call trace_and_save_fn,$(client_agent_port),heap,client_heap.out)


server_run:
	# go run cmd/server/main.go
	go run -ldflags "-X main.buildVersion=$$(cat cmd/server/version.txt) -X 'main.buildDate=$$(date +'%d/%m/%Y')' -X 'main.buildCommit=$$(git rev-parse HEAD)'" cmd/server/main.go
server_trace_profile:
	$(call trace_fn,$(server_agent_port),profile)
server_trace_profile_and_save:
	$(call trace_and_save_fn,$(server_agent_port),profile,server_profile.out)
server_trace_heap:
	$(call trace_fn,$(server_agent_port),heap)
server_trace_heap_and_save:
	$(call trace_and_save_fn,$(server_agent_port),heap,server_heap.out)

run_doc:
	(sleep 2; open "http://localhost:6060/pkg/?m=all")&
	~/go/bin/godoc  -http=localhost:6060 -goroot=. -play

run_lint:
	go run ./multichecker.go /Users/a.bogushev/course/yametrics/...
