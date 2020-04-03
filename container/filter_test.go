package container

import (
	"regexp"
	"testing"

	"docker.io/go-docker/api/types/events"
	"github.com/stretchr/testify/assert"
)

func Test_newFilter(t *testing.T) {

	cfg := DockerCfg{
		ContainerEvents:  false,
		VolumeEvents:     false,
		NetworkEvents:    false,
		ShowRunning:      false,
		ContainerActions: []string{"test", "test1"},
		NetworkActions:   []string{"test2", "test3"},
		VolumeActions:    []string{"test4", "test5"},
	}
	f, err := newFilter(cfg)
	assert.Nil(t, err)
	assert.NotNil(t, f)
	assert.Equal(t, 6, len(f.actions))

	_, ok := f.actions[containerEvent+"test"]
	assert.True(t, ok)
	_, ok = f.actions[containerEvent+"test1"]
	assert.True(t, ok)

	_, ok = f.actions[networkEvent+"test2"]
	assert.True(t, ok)
	_, ok = f.actions[networkEvent+"test3"]
	assert.True(t, ok)

	_, ok = f.actions[volumeEvent+"test4"]
	assert.True(t, ok)
	_, ok = f.actions[volumeEvent+"test5"]
	assert.True(t, ok)
}

func Test_newFilterBadRegExp(t *testing.T) {
	cfg := DockerCfg{

		FilterName: `^\/(?!\/)(.*?)`,
	}

	_, err := newFilter(cfg)
	assert.Error(t, err)

	cfg = DockerCfg{

		FilterImage: `^\/(?!\/)(.*?)`,
	}

	_, err = newFilter(cfg)
	assert.Error(t, err)
}

func Test_filter_accept(t *testing.T) {
	type fields struct {
		actions       map[string]struct{}
		imageName     *regexp.Regexp
		containerName *regexp.Regexp
	}

	type args struct {
		event events.Message
	}

	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "Event on empty filter",
			fields: fields{
				actions: map[string]struct{}{},
			},
			args: args{event: events.Message{

				Type:   "container",
				Action: "kill",
			},
			},
			want: false,
		},
		{
			name: "Event on action not present filter",
			fields: fields{
				actions: map[string]struct{}{"container_start": {}},
			},
			args: args{event: events.Message{

				Type:   "container",
				Action: "kill",
			},
			},
			want: false,
		},
		{
			name: "Event on action  present ",
			fields: fields{
				actions: map[string]struct{}{"container_start": {}},
			},
			args: args{event: events.Message{

				Type:   "container",
				Action: "start",
			},
			},
			want: true,
		},
		{
			name: "Event on action  present 2",
			fields: fields{
				actions: map[string]struct{}{"network_create": {}},
			},
			args: args{event: events.Message{

				Type:   "network",
				Action: "create",
			},
			},
			want: true,
		},
		{
			name: "Event on action  present 2",
			fields: fields{
				actions: map[string]struct{}{"network_create": {}},
			},
			args: args{event: events.Message{

				Type:   "network",
				Action: "create",
			},
			},
			want: true,
		},
		{
			name: "Event on action  exec 2",
			fields: fields{
				actions: map[string]struct{}{"container_exec_start": {}},
			},
			args: args{event: events.Message{

				Type:   "container",
				Action: "exec_start: bash",
			},
			},
			want: true,
		},
		{
			name: "Filter by events and image name 1",
			fields: fields{
				actions:       map[string]struct{}{"container_exec_start": {}},
				imageName:     regexp.MustCompile(`^hello.?world$`),
				containerName: nil,
			},
			args: args{event: events.Message{

				Type:   "container",
				Action: "exec_start: bash",
				Actor: events.Actor{
					Attributes: map[string]string{"image": "hello-world"},
				},
			},
			},
			want: true,
		},
		{
			name: "Filter by events and image name 2",
			fields: fields{
				actions:       map[string]struct{}{"container_exec_start": {}},
				imageName:     regexp.MustCompile(`^hello-world$`),
				containerName: nil,
			},
			args: args{event: events.Message{

				Type:   "container",
				Action: "exec_start: bash",
				Actor: events.Actor{
					Attributes: map[string]string{"image": "Myhello-world"},
				},
			},
			},
			want: false,
		},
		{
			name: "Filter by events and container name 1",
			fields: fields{
				actions:       map[string]struct{}{"container_exec_start": {}},
				imageName:     regexp.MustCompile(`^hello-world_OK$`),
				containerName: nil,
			},
			args: args{event: events.Message{

				Type:   "container",
				Action: "exec_start: bash",
				Actor: events.Actor{
					Attributes: map[string]string{"image": "hello-world_OK"},
				},
			},
			},
			want: true,
		},
		{
			name: "Filter by events and container name 2",
			fields: fields{
				actions:       map[string]struct{}{"container_exec_start": {}},
				imageName:     regexp.MustCompile(`^hello-world_OK$`),
				containerName: nil,
			},
			args: args{event: events.Message{

				Type:   "container",
				Action: "exec_start: bash",
				Actor: events.Actor{
					Attributes: map[string]string{"image": "NOhello-world_OK"},
				},
			},
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &filter{
				actions:       tt.fields.actions,
				imageName:     tt.fields.imageName,
				containerName: tt.fields.containerName,
			}
			if got := f.accept(tt.args.event); got != tt.want {
				t.Errorf("accept() = %v, want %v", got, tt.want)
			}
		})
	}
}
