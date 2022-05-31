// +build ignore

//#include <linux/kernel.h>
#include "common.h"

char __license[] SEC("license") = "Dual MIT/GPL";

struct bpf_map_def SEC("maps") kp_array_map = {
	.type        = BPF_MAP_TYPE_ARRAY,
	.key_size    = sizeof(u32),
	.value_size  = sizeof(u64),
	.max_entries = 100,
};

//char* format = "<5> bpf_map_update_elem executed !\n";

SEC("kprobe/sys_execve")
int kprobe_sys_execve() {
	u32 key     = 0;
	u64 initval = 1, *valp;

	valp = bpf_map_lookup_elem(&kp_array_map, &key);
	if (!valp) {
		bpf_map_update_elem(&kp_array_map, &key, &initval, BPF_ANY);
//		bpf_trace_printk(format, sizeof(format));
        // map .rodata: map create: read- and write-only maps not supported (requires >= v5.2)
        // bpf_printk("<5> bpf_map_update_elem executed !\n");
		return 0;
	}
	__sync_fetch_and_add(valp, 1);

	return 0;
}

SEC("kprobe/sys_connect")
int kprobe_sys_connect() {
	u32 key     = 1;
	u64 initval = 1, *valp;

	valp = bpf_map_lookup_elem(&kp_array_map, &key);
	if (!valp) {
		bpf_map_update_elem(&kp_array_map, &key, &initval, BPF_ANY);
//		bpf_trace_printk(format, sizeof(format));
        // map .rodata: map create: read- and write-only maps not supported (requires >= v5.2)
        // bpf_printk("<5> bpf_map_update_elem executed !\n");
		return 0;
	}
	__sync_fetch_and_add(valp, 1);

	return 0;
}

SEC("kprobe/sys_accept")
int kprobe_sys_accept() {
	u32 key     = 2;
	u64 initval = 1, *valp;

	valp = bpf_map_lookup_elem(&kp_array_map, &key);
	if (!valp) {
		bpf_map_update_elem(&kp_array_map, &key, &initval, BPF_ANY);
//		bpf_trace_printk(format, sizeof(format));
        // map .rodata: map create: read- and write-only maps not supported (requires >= v5.2)
        // bpf_printk("<5> bpf_map_update_elem executed !\n");
		return 0;
	}
	__sync_fetch_and_add(valp, 1);

	return 0;
}

#define MAX_LENGTH 128

struct msg {
	__s32 seq;
	__u64 cts;
	__u8 comm[MAX_LENGTH];
};

struct bpf_map_def SEC("maps") kp_perf_events = {
	.type = BPF_MAP_TYPE_PERF_EVENT_ARRAY,
	.key_size = sizeof(u32),
	.value_size = sizeof(u32),
	.max_entries = 128,
};

SEC("kprobe/vfs_read")
int kprobe_vfs_read(struct pt_regs *ctx) {
    // program array
	u32 key     = 3;
    u64 initval = 1, *valp;
    valp = bpf_map_lookup_elem(&kp_array_map, &key);
    if (!valp) {
    	bpf_map_update_elem(&kp_array_map, &key, &initval, BPF_ANY);
    	return 0;
    }
    __sync_fetch_and_add(valp, 1);

//   char msg[] = "hello BPF!\n";
//   bpf_trace_printk(msg, sizeof(msg));

    // perf event array
//	unsigned long cts = bpf_ktime_get_ns();
//	struct msg val = {0};
//	static __u32 seq = 0;
//	val.seq = seq = (seq + 1) % 4294967295U;
//	val.cts = bpf_ktime_get_ns();
//	bpf_get_current_comm(val.comm, sizeof(val.comm));
//	bpf_perf_event_output(ctx, &kp_perf_events, 0, &val, sizeof(val));

	return 0;
}
