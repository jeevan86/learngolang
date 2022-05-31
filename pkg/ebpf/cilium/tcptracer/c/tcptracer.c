#undef  __i386__

#include "kconfig.h"    // #define IS_ENABLED ...
#include "common.h"     // struct pt_regs ...
//#include <stdbool.h>  // bool  TRUE=1
#include <stddef.h>
#include <stdint.h>
//#include <linux/kconfig.h>

// https://zhuanlan.zhihu.com/p/39869242
// https://blog.csdn.net/weibo1230123/article/details/84106028
//#define __packed
//
//typedef struct {
//        int counter;
//} atomic_t;
//
typedef struct {
        long counter;
} atomic64_t;

#define __printf(a, b)          __attribute__((format(printf, a, b)))
#define __must_check            __attribute__((warn_unused_result))
#define __malloc                __attribute__((__malloc__))

#pragma clang diagnostic push
#pragma clang diagnostic ignored "-Wgnu-variable-sized-type-not-at-end"
#pragma clang diagnostic ignored "-Waddress-of-packed-member"
//#include <linux/ptrace.h>
#pragma clang diagnostic pop
//#include <linux/version.h>
#include <linux/types.h>
//#include <linux/bpf.h>
#pragma clang diagnostic push
#pragma clang diagnostic ignored "-Wtautological-compare"
#pragma clang diagnostic ignored "-Wgnu-variable-sized-type-not-at-end"
#pragma clang diagnostic ignored "-Wenum-conversion"
//#include <net/sock.h>
#pragma clang diagnostic pop
//#include <net/inet_sock.h>
//#include <net/net_namespace.h>
//#include <linux/spinlock.h>
#include "spinlock.h"
#include "wait.h"

#include "tcptracer.h"

/* This is a key/value store with the keys being the cpu number
 * and the values being a perf file descriptor.
 */
 struct bpf_map_def SEC("maps/tcp_event_ipv4") tcp_event_ipv4 = {
	.type = BPF_MAP_TYPE_PERF_EVENT_ARRAY,
	.key_size = sizeof(int),
	.value_size = sizeof(__u32),
	.max_entries = 1024,
	.pinning = 0,
	.namespace = "",
};

/* This is a key/value store with the keys being the cpu number
 * and the values being a perf file descriptor.
 */
struct bpf_map_def SEC("maps/tcp_event_ipv6") tcp_event_ipv6 = {
	.type = BPF_MAP_TYPE_PERF_EVENT_ARRAY,
	.key_size = sizeof(int),
	.value_size = sizeof(__u32),
	.max_entries = 1024,
	.pinning = 0,
	.namespace = "",
};

/* This is a key/value store with the keys being an ipv4_tuple_t
 * and the values being a struct pid_comm_t.
 */
struct bpf_map_def SEC("maps/tuplepid_ipv4") tuplepid_ipv4 = {
	.type = BPF_MAP_TYPE_HASH,
	.key_size = sizeof(struct ipv4_tuple_t),
	.value_size = sizeof(struct pid_comm_t),
	.max_entries = 1024,
	.pinning = 0,
	.namespace = "",
};

/* This is a key/value store with the keys being an ipv6_tuple_t
 * and the values being a struct pid_comm_t.
 */
struct bpf_map_def SEC("maps/tuplepid_ipv6") tuplepid_ipv6 = {
	.type = BPF_MAP_TYPE_HASH,
	.key_size = sizeof(struct ipv6_tuple_t),
	.value_size = sizeof(struct pid_comm_t),
	.max_entries = 1024,
	.pinning = 0,
	.namespace = "",
};

struct bpf_map_def SEC("maps/tcptracer_status") tcptracer_status = {
	.type = BPF_MAP_TYPE_HASH,
	.key_size = sizeof(__u64),
	.value_size = sizeof(struct tcptracer_status_t),
	.max_entries = 1,
	.pinning = 0,
	.namespace = "",
};

SEC("kprobe/tcp_set_state")
int tcptracer_tcp_set_state(struct pt_regs *ctx)
{

//   char msg[] = "hello BPF!\n";
//   bpf_trace_printk(msg, sizeof(msg));

	u32 cpu = bpf_get_smp_processor_id();
	struct bpf_sock *skp; // struct sock *skp; // #include <net/sock.h>
	struct tcptracer_status_t *status;     // custom header file
	int state;
	u64 zero = 0;
	skp =  (struct bpf_sock *) PT_REGS_PARM1(ctx); // bpf_helpers.h
	state = (int) PT_REGS_PARM2(ctx);

	status = bpf_map_lookup_elem(&tcptracer_status, &zero);
	if (status == NULL || status->state != TCPTRACER_STATE_READY) {
		return 0;
	}

	if (state != TCP_ESTABLISHED && state != TCP_CLOSE) {
		return 0;
	}

	if (check_family(skp, AF_INET)) {
		// output
		struct ipv4_tuple_t t = { };
		if (!read_ipv4_tuple(&t, status, skp)) {
			return 0;
		}
		if (state == TCP_CLOSE) {
			bpf_map_delete_elem(&tuplepid_ipv4, &t);
			return 0;
		}

		struct pid_comm_t *pp;

		pp = bpf_map_lookup_elem(&tuplepid_ipv4, &t);
		if (pp == 0) {
			return 0;	// missed entry
		}
		struct pid_comm_t p = { };
		bpf_probe_read(&p, sizeof(struct pid_comm_t), pp);

		struct tcp_ipv4_event_t evt4 = {
			.timestamp = bpf_ktime_get_ns(),
			.cpu = cpu,
			.type = TCP_EVENT_TYPE_CONNECT,
			.pid = p.pid >> 32,
			.saddr = t.saddr,
			.daddr = t.daddr,
			.sport = ntohs(t.sport),
			.dport = ntohs(t.dport),
			.netns = t.netns,
		};
		int i;
		for (i = 0; i < TASK_COMM_LEN; i++) {
			evt4.comm[i] = p.comm[i];
		}

		bpf_perf_event_output(ctx, &tcp_event_ipv4, cpu, &evt4, sizeof(evt4));
		bpf_map_delete_elem(&tuplepid_ipv4, &t);
	} else if (check_family(skp, AF_INET6)) {
		// output
		struct ipv6_tuple_t t = { };
		if (!read_ipv6_tuple(&t, status, skp)) {
			return 0;
		}
		if (state == TCP_CLOSE) {
			bpf_map_delete_elem(&tuplepid_ipv6, &t);
			return 0;
		}

		struct pid_comm_t *pp;
		pp = bpf_map_lookup_elem(&tuplepid_ipv6, &t);
		if (pp == 0) {
			return 0;       // missed entry
		}
		struct pid_comm_t p = { };
		bpf_probe_read(&p, sizeof(struct pid_comm_t), pp);
		struct tcp_ipv6_event_t evt6 = {
			.timestamp = bpf_ktime_get_ns(),
			.cpu = cpu,
			.type = TCP_EVENT_TYPE_CONNECT,
			.pid = p.pid >> 32,
			.saddr_h = t.saddr_h,
			.saddr_l = t.saddr_l,
			.daddr_h = t.daddr_h,
			.daddr_l = t.daddr_l,
			.sport = ntohs(t.sport),
			.dport = ntohs(t.dport),
			.netns = t.netns,
		};
		int i;
		for (i = 0; i < TASK_COMM_LEN; i++) {
			evt6.comm[i] = p.comm[i];
		}

		bpf_perf_event_output(ctx, &tcp_event_ipv6, cpu, &evt6, sizeof(evt6));
		bpf_map_delete_elem(&tuplepid_ipv6, &t);
	}

	return 0;
}

__attribute__((always_inline))
static int read_ipv4_tuple(struct ipv4_tuple_t *tuple, struct tcptracer_status_t *status, struct sock *skp)
{
	u32 saddr, daddr, net_ns_inum;
	u16 sport, dport;
	possible_net_t *skc_net;

	saddr = 0;
	daddr = 0;
	sport = 0;
	dport = 0;
	skc_net = NULL;
	net_ns_inum = 0;

	bpf_probe_read(&saddr, sizeof(saddr), ((char *)skp) + status->offset_saddr);
	bpf_probe_read(&daddr, sizeof(daddr), ((char *)skp) + status->offset_daddr);
	bpf_probe_read(&sport, sizeof(sport), ((char *)skp) + status->offset_sport);
	bpf_probe_read(&dport, sizeof(dport), ((char *)skp) + status->offset_dport);
	// Get network namespace id
	bpf_probe_read(&skc_net, sizeof(void *), ((char *)skp) + status->offset_netns);
	bpf_probe_read(&net_ns_inum, sizeof(net_ns_inum), ((char *)skc_net) + status->offset_ino);

	tuple->saddr = saddr;
	tuple->daddr = daddr;
	tuple->sport = sport;
	tuple->dport = dport;
	tuple->netns = net_ns_inum;

	// if addresses or ports are 0, ignore
	if (saddr == 0 || daddr == 0 || sport == 0 || dport == 0) {
		return 0;
	}

	return 1;
}

__attribute__((always_inline))
static int read_ipv6_tuple(struct ipv6_tuple_t *tuple, struct tcptracer_status_t *status, struct sock *skp)
{
	u32 net_ns_inum;
	u16 sport, dport;
	u64 saddr_h, saddr_l, daddr_h, daddr_l;
	possible_net_t *skc_net;

	saddr_h = 0;
	saddr_l = 0;
	daddr_h = 0;
	daddr_l = 0;
	sport = 0;
	dport = 0;
	skc_net = NULL;
	net_ns_inum = 0;

	bpf_probe_read(&saddr_h, sizeof(saddr_h), ((char *)skp) + status->offset_daddr_ipv6 + 2 * sizeof(u64));
	bpf_probe_read(&saddr_l, sizeof(saddr_l), ((char *)skp) + status->offset_daddr_ipv6 + 3 * sizeof(u64));
	bpf_probe_read(&daddr_h, sizeof(daddr_h), ((char *)skp) + status->offset_daddr_ipv6);
	bpf_probe_read(&daddr_l, sizeof(daddr_l), ((char *)skp) + status->offset_daddr_ipv6 + sizeof(u64));
	bpf_probe_read(&sport, sizeof(sport), ((char *)skp) + status->offset_sport);
	bpf_probe_read(&dport, sizeof(dport), ((char *)skp) + status->offset_dport);
	// Get network namespace id
	bpf_probe_read(&skc_net, sizeof(void *), ((char *)skp) + status->offset_netns);
	bpf_probe_read(&net_ns_inum, sizeof(net_ns_inum), ((char *)skc_net) + status->offset_ino);

	tuple->saddr_h = saddr_h;
	tuple->saddr_l = saddr_l;
	tuple->daddr_h = daddr_h;
	tuple->daddr_l = daddr_l;
	tuple->sport = sport;
	tuple->dport = dport;
	tuple->netns = net_ns_inum;

	// if addresses or ports are 0, ignore
	if (!(saddr_h || saddr_l) || !(daddr_h || daddr_l) || sport == 0 || dport == 0) {
		return 0;
	}

	return 1;
}