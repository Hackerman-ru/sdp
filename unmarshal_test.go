// SPDX-FileCopyrightText: 2023 The Pion community <https://pion.ly>
// SPDX-License-Identifier: MIT

package sdp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	BaseSDP = "v=0\r\n" +
		"o=jdoe 2890844526 2890842807 IN IP4 10.47.16.5\r\n" +
		"s=SDP Seminar\r\n"

	SessionInformationSDP = BaseSDP +
		"i=A Seminar on the session description protocol\r\n" +
		"t=3034423619 3042462419\r\n"

	// https://tools.ietf.org/html/rfc4566#section-5
	// Parsers SHOULD be tolerant and also accept records terminated
	// with a single newline character.
	SessionInformationSDPLFOnly = "v=0\n" +
		"o=jdoe 2890844526 2890842807 IN IP4 10.47.16.5\n" +
		"s=SDP Seminar\n" +
		"i=A Seminar on the session description protocol\n" +
		"t=3034423619 3042462419\n"

	// SessionInformationSDPCROnly = "v=0\r" +
	// 	"o=jdoe 2890844526 2890842807 IN IP4 10.47.16.5\r" +
	// 	"s=SDP Seminar\r"
	// 	"i=A Seminar on the session description protocol\r" +
	// 	"t=3034423619 3042462419\r"

	// Other SDP parsers (e.g. one in VLC media player) allow
	// empty lines.
	SessionInformationSDPExtraCRLF = "v=0\r\n" +
		"o=jdoe 2890844526 2890842807 IN IP4 10.47.16.5\r\n" +
		"\r\n" +
		"s=SDP Seminar\r\n" +
		"\r\n" +
		"i=A Seminar on the session description protocol\r\n" +
		"\r\n" +
		"t=3034423619 3042462419\r\n" +
		"\r\n"

	URISDP = BaseSDP +
		"u=http://www.example.com/seminars/sdp.pdf\r\n" +
		"t=3034423619 3042462419\r\n"

	EmailAddressSDP = BaseSDP +
		"e=j.doe@example.com (Jane Doe)\r\n" +
		"t=3034423619 3042462419\r\n"

	PhoneNumberSDP = BaseSDP +
		"p=+1 617 555-6011\r\n" +
		"t=3034423619 3042462419\r\n"

	SessionConnectionInformationSDP = BaseSDP +
		"c=IN IP4 224.2.17.12/127\r\n" +
		"t=3034423619 3042462419\r\n"

	SessionBandwidthSDP = BaseSDP +
		"b=X-YZ:128\r\n" +
		"b=AS:12345\r\n" +
		"t=3034423619 3042462419\r\n"

	TimingSDP = BaseSDP +
		"t=2873397496 2873404696\r\n"

	// Short hand time notation is converted into NTP timestamp format in
	// seconds. Because of that unittest comparisons will fail as the same time
	// will be expressed in different units.
	RepeatTimesSDP = TimingSDP +
		"r=604800 3600 0 90000\r\n" +
		"r=3d 2h 0 21h\r\n"

	RepeatTimesSDPExpected = TimingSDP +
		"r=604800 3600 0 90000\r\n" +
		"r=259200 7200 0 75600\r\n"

	RepeatTimesSDPExtraCRLF = RepeatTimesSDPExpected +
		"\r\n"

	// The expected value looks a bit different for the same reason as mentioned
	// above regarding RepeatTimes.
	TimeZonesSDP = TimingSDP +
		"r=2882844526 -1h 2898848070 0\r\n"

	TimeZonesSDPExpected = TimingSDP +
		"r=2882844526 -3600 2898848070 0\r\n"

	TimeZonesSDP2 = TimingSDP +
		"z=2882844526 -3600 2898848070 0\r\n"

	TimeZonesSDP2ExtraCRLF = TimeZonesSDP2 +
		"\r\n"

	SessionEncryptionKeySDP = TimingSDP +
		"k=prompt\r\n"

	SessionEncryptionKeySDPExtraCRLF = SessionEncryptionKeySDP +
		"\r\n"

	SessionAttributesSDP = TimingSDP +
		"a=rtpmap:96 opus/48000\r\n"

	MediaNameSDP = TimingSDP +
		"m=video 51372 RTP/AVP 99\r\n" +
		"m=audio 54400 RTP/SAVPF 0 96\r\n" +
		"m=message 5028 TCP/MSRP *\r\n"

	MediaNameSDPExtraCRLF = MediaNameSDP +
		"\r\n"

	MediaTitleSDP = MediaNameSDP +
		"i=Vivamus a posuere nisl\r\n"

	MediaConnectionInformationSDP = MediaNameSDP +
		"c=IN IP4 203.0.113.1\r\n"

	MediaConnectionInformationSDPExtraCRLF = MediaConnectionInformationSDP +
		"\r\n"

	MediaDescriptionOutOfOrderSDP = MediaNameSDP +
		"a=rtpmap:99 h263-1998/90000\r\n" +
		"a=candidate:0 1 UDP 2113667327 203.0.113.1 54400 typ host\r\n" +
		"c=IN IP4 203.0.113.1\r\n" +
		"i=Vivamus a posuere nisl\r\n"

	MediaDescriptionOutOfOrderSDPActual = MediaNameSDP +
		"i=Vivamus a posuere nisl\r\n" +
		"c=IN IP4 203.0.113.1\r\n" +
		"a=rtpmap:99 h263-1998/90000\r\n" +
		"a=candidate:0 1 UDP 2113667327 203.0.113.1 54400 typ host\r\n"

	MediaBandwidthSDP = MediaNameSDP +
		"b=X-YZ:128\r\n" +
		"b=AS:12345\r\n" +
		"b=TIAS:12345\r\n" +
		"b=RS:12345\r\n" +
		"b=RR:12345\r\n"

	MediaEncryptionKeySDP = MediaNameSDP +
		"k=prompt\r\n"

	MediaEncryptionKeySDPExtraCRLF = MediaEncryptionKeySDP +
		"\r\n"

	MediaAttributesSDP = MediaNameSDP +
		"a=rtpmap:99 h263-1998/90000\r\n" +
		"a=candidate:0 1 UDP 2113667327 203.0.113.1 54400 typ host\r\n" +
		"a=rtcp-fb:97 ccm fir\r\n" +
		"a=rtcp-fb:97 nack\r\n" +
		"a=rtcp-fb:97 nack pli\r\n"

	MediaBfcpSDP = TimingSDP +
		"m=application 3238 UDP/BFCP *\r\n" +
		"a=sendrecv\r\n" +
		"a=setup:actpass\r\n" +
		"a=connection:new\r\n" +
		"a=floorctrl:c-s\r\n"

	MediaCubeSDP = TimingSDP +
		"m=application 2455 UDP/UDT/IX *\r\n" +
		"a=ixmap:0 ping\r\n" +
		"a=ixmap:2 xccp\r\n"

	MediaTCPMRCPv2 = TimingSDP +
		"m=application 1544 TCP/MRCPv2 1\r\n"

	MediaTCPTLSMRCPv2 = TimingSDP +
		"m=application 1544 TCP/TLS/MRCPv2 1\r\n"

	CanonicalUnmarshalSDP = "v=0\r\n" +
		"o=jdoe 2890844526 2890842807 IN IP4 10.47.16.5\r\n" +
		"s=SDP Seminar\r\n" +
		"i=A Seminar on the session description protocol\r\n" +
		"u=http://www.example.com/seminars/sdp.pdf\r\n" +
		"e=j.doe@example.com (Jane Doe)\r\n" +
		"p=+1 617 555-6011\r\n" +
		"c=IN IP4 224.2.17.12/127\r\n" +
		"b=X-YZ:128\r\n" +
		"b=AS:12345\r\n" +
		"t=2873397496 2873404696\r\n" +
		"t=3034423619 3042462419\r\n" +
		"r=604800 3600 0 90000\r\n" +
		"z=2882844526 -3600 2898848070 0\r\n" +
		"k=prompt\r\n" +
		"a=candidate:0 1 UDP 2113667327 203.0.113.1 54400 typ host\r\n" +
		"a=recvonly\r\n" +
		"m=audio 49170 RTP/AVP 0\r\n" +
		"i=Vivamus a posuere nisl\r\n" +
		"c=IN IP4 203.0.113.1\r\n" +
		"b=X-YZ:128\r\n" +
		"k=prompt\r\n" +
		"a=sendrecv\r\n" +
		"m=video 51372 RTP/AVP 99\r\n" +
		"a=rtpmap:99 h263-1998/90000\r\n"
)

func TestRoundTrip(t *testing.T) {
	for _, test := range []struct {
		Name   string
		SDP    string
		Actual string
	}{
		{
			Name:   "SessionInformationSDPLFOnly",
			SDP:    SessionInformationSDPLFOnly,
			Actual: SessionInformationSDP,
		},
		// {
		// 	Name:   "SessionInformationSDPCROnly",
		// 	SDP:    SessionInformationSDPCROnly,
		// 	Actual: SessionInformationSDPBaseSDP,
		// },
		{
			Name:   "SessionInformationSDPExtraCRLF",
			SDP:    SessionInformationSDPExtraCRLF,
			Actual: SessionInformationSDP,
		},
		{
			Name: "SessionInformation",
			SDP:  SessionInformationSDP,
		},
		{
			Name: "URI",
			SDP:  URISDP,
		},
		{
			Name: "EmailAddress",
			SDP:  EmailAddressSDP,
		},
		{
			Name: "PhoneNumber",
			SDP:  PhoneNumberSDP,
		},
		{
			Name:   "RepeatTimesSDPExtraCRLF",
			SDP:    RepeatTimesSDPExtraCRLF,
			Actual: RepeatTimesSDPExpected,
		},
		{
			Name: "SessionConnectionInformation",
			SDP:  SessionConnectionInformationSDP,
		},
		{
			Name: "SessionBandwidth",
			SDP:  SessionBandwidthSDP,
		},
		{
			Name: "SessionEncryptionKey",
			SDP:  SessionEncryptionKeySDP,
		},
		{
			Name:   "SessionEncryptionKeyExtraCRLF",
			SDP:    SessionEncryptionKeySDPExtraCRLF,
			Actual: SessionEncryptionKeySDP,
		},
		{
			Name: "SessionAttributes",
			SDP:  SessionAttributesSDP,
		},
		{
			Name:   "TimeZonesSDP2ExtraCRLF",
			SDP:    TimeZonesSDP2ExtraCRLF,
			Actual: TimeZonesSDP2,
		},
		{
			Name: "MediaName",
			SDP:  MediaNameSDP,
		},
		{
			Name:   "MediaNameExtraCRLF",
			SDP:    MediaNameSDPExtraCRLF,
			Actual: MediaNameSDP,
		},
		{
			Name: "MediaTitle",
			SDP:  MediaTitleSDP,
		},
		{
			Name: "MediaConnectionInformation",
			SDP:  MediaConnectionInformationSDP,
		},
		{
			Name:   "MediaConnectionInformationExtraCRLF",
			SDP:    MediaConnectionInformationSDPExtraCRLF,
			Actual: MediaConnectionInformationSDP,
		},
		{
			Name:   "MediaDescriptionOutOfOrder",
			SDP:    MediaDescriptionOutOfOrderSDP,
			Actual: MediaDescriptionOutOfOrderSDPActual,
		},
		{
			Name: "MediaBandwidth",
			SDP:  MediaBandwidthSDP,
		},
		{
			Name: "MediaEncryptionKey",
			SDP:  MediaEncryptionKeySDP,
		},
		{
			Name:   "MediaEncryptionKeyExtraCRLF",
			SDP:    MediaEncryptionKeySDPExtraCRLF,
			Actual: MediaEncryptionKeySDP,
		},
		{
			Name: "MediaAttributes",
			SDP:  MediaAttributesSDP,
		},
		{
			Name: "CanonicalUnmarshal",
			SDP:  CanonicalUnmarshalSDP,
		},
		{
			Name: "MediaBfcpSDP",
			SDP:  MediaBfcpSDP,
		},
		{
			Name: "MediaCubeSDP",
			SDP:  MediaCubeSDP,
		},
		{
			Name: "MediaTCPMRCPv2",
			SDP:  MediaTCPMRCPv2,
		},
		{
			Name: "MediaTCPTLSMRCPv2",
			SDP:  MediaTCPTLSMRCPv2,
		},
	} {
		test := test
		t.Run(test.Name, func(t *testing.T) {
			sd := &SessionDescription{}

			err := sd.UnmarshalString(test.SDP)
			assert.NoError(t, err)

			actual, err := sd.Marshal()
			assert.NoError(t, err)

			want := test.SDP
			if test.Actual != "" {
				want = test.Actual
			}

			assert.Equal(t, want, string(actual))
		})
	}
}

func TestUnmarshalRepeatTimes(t *testing.T) {
	sd := &SessionDescription{}
	assert.NoError(t, sd.UnmarshalString(RepeatTimesSDP))

	actual, err := sd.Marshal()
	assert.NoError(t, err)
	assert.Equal(t, RepeatTimesSDPExpected, string(actual))

	err = sd.UnmarshalString(TimingSDP + "r=\r\n")
	assert.ErrorIs(t, err, errSDPInvalidValue)
}

func TestUnmarshalTimeZones(t *testing.T) {
	sd := &SessionDescription{}
	assert.NoError(t, sd.UnmarshalString(TimeZonesSDP))

	actual, err := sd.Marshal()
	assert.NoError(t, err)
	assert.Equal(t, TimeZonesSDPExpected, string(actual))
}

func TestUnmarshalNonNilAddress(t *testing.T) {
	in := "v=0\r\no=0 0 0 IN IP4 0\r\ns=0\r\nc=IN IP4\r\nt=0 0\r\n"
	var sd SessionDescription
	err := sd.UnmarshalString(in)
	assert.NoError(t, err)

	out, err := sd.Marshal()
	assert.NoError(t, err)
	assert.Equal(t, in, string(out))
}

func TestUnmarshalZeroValues(t *testing.T) {
	in := "v=0\r\no=0 0 0 IN IP4 0\r\ns=\r\nt=0 0\r\n"
	var sd SessionDescription
	assert.NoError(t, sd.UnmarshalString(in))

	out, err := sd.Marshal()
	assert.NoError(t, err)
	assert.Equal(t, in, string(out))
}

func TestUnmarshalPortRange(t *testing.T) {
	for _, test := range []struct {
		In          string
		ExpectError error
	}{
		{
			In:          SessionAttributesSDP + "m=video -1 RTP/AVP 99\r\n",
			ExpectError: errSDPInvalidPortValue,
		},
		{
			In:          SessionAttributesSDP + "m=video 65536 RTP/AVP 99\r\n",
			ExpectError: errSDPInvalidPortValue,
		},
		{
			In:          SessionAttributesSDP + "m=video 0 RTP/AVP 99\r\n",
			ExpectError: nil,
		},
		{
			In:          SessionAttributesSDP + "m=video 65535 RTP/AVP 99\r\n",
			ExpectError: nil,
		},
		{
			In:          SessionAttributesSDP + "m=video --- RTP/AVP 99\r\n",
			ExpectError: errSDPInvalidPortValue,
		},
	} {
		var sd SessionDescription
		err := sd.UnmarshalString(test.In)
		if test.ExpectError != nil {
			assert.ErrorIs(t, err, test.ExpectError)
		} else {
			assert.NoError(t, err)
		}
	}
}

func BenchmarkUnmarshal(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		var sd SessionDescription
		err := sd.UnmarshalString(CanonicalUnmarshalSDP)
		assert.NoError(b, err)
	}
}

func TestUnmarshalOriginIncomplete(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected Origin
	}{
		{
			name:  "missing unicast address - Uniview camera case",
			input: "v=0\r\no=- 1001 1 IN IP4\r\ns=VCP IPC Realtime stream\r\nt=0 0\r\n",
			expected: Origin{
				Username:       "-",
				SessionID:      1001,
				SessionVersion: 1,
				NetworkType:    "IN",
				AddressType:    "IP4",
				UnicastAddress: "0.0.0.0",
			},
		},
		{
			name:  "missing address type and address",
			input: "v=0\r\no=- 1001 1 IN\r\ns=Test Stream\r\nt=0 0\r\n",
			expected: Origin{
				Username:       "-",
				SessionID:      1001,
				SessionVersion: 1,
				NetworkType:    "IN",
				AddressType:    "IP4",
				UnicastAddress: "0.0.0.0",
			},
		},
		{
			name:  "IPv6 missing address",
			input: "v=0\r\no=- 1001 1 IN IP6\r\ns=Test Stream\r\nt=0 0\r\n",
			expected: Origin{
				Username:       "-",
				SessionID:      1001,
				SessionVersion: 1,
				NetworkType:    "IN",
				AddressType:    "IP6",
				UnicastAddress: "::",
			},
		},
		{
			name:  "complete origin line - should work as before",
			input: "v=0\r\no=jdoe 2890844526 2890842807 IN IP4 10.47.16.5\r\ns=SDP Seminar\r\nt=3034423619 3042462419\r\n",
			expected: Origin{
				Username:       "jdoe",
				SessionID:      2890844526,
				SessionVersion: 2890842807,
				NetworkType:    "IN",
				AddressType:    "IP4",
				UnicastAddress: "10.47.16.5",
			},
		},
		{
			name:  "empty address field",
			input: "v=0\r\no=- 1001 1 IN IP4 \r\ns=Test\r\nt=0 0\r\n",
			expected: Origin{
				Username:       "-",
				SessionID:      1001,
				SessionVersion: 1,
				NetworkType:    "IN",
				AddressType:    "IP4",
				UnicastAddress: "0.0.0.0",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var sd SessionDescription
			err := sd.UnmarshalString(test.input)
			assert.NoError(t, err)
			assert.Equal(t, test.expected, sd.Origin)
		})
	}
}

func TestUnmarshalOriginInvalidFields(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "invalid network type",
			input: "v=0\r\no=- 1001 1 INVALID IP4 10.0.0.1\r\ns=Test\r\nt=0 0\r\n",
		},
		{
			name:  "invalid address type",
			input: "v=0\r\no=- 1001 1 IN INVALID 10.0.0.1\r\ns=Test\r\nt=0 0\r\n",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var sd SessionDescription
			err := sd.UnmarshalString(test.input)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "invalid value")
		})
	}
}

// Test edge cases for 100% coverage.
func TestUnmarshalOriginEdgeCases(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectError bool
	}{
		{
			name:        "missing mandatory username",
			input:       "v=0\r\no=\r\ns=Test\r\nt=0 0\r\n",
			expectError: true,
		},
		{
			name:        "missing mandatory session ID",
			input:       "v=0\r\no=user\r\ns=Test\r\nt=0 0\r\n",
			expectError: true,
		},
		{
			name:        "missing mandatory network type",
			input:       "v=0\r\no=user 1001 1\r\ns=Test\r\nt=0 0\r\n",
			expectError: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var sd SessionDescription
			err := sd.UnmarshalString(test.input)
			if test.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
