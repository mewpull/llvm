// Code generated by "stringer -linecomment -type CallingConv"; DO NOT EDIT.

package enum

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[CallingConvNone-0]
	_ = x[CallingConvC-1]
	_ = x[CallingConvFast-8]
	_ = x[CallingConvCold-9]
	_ = x[CallingConvGHC-10]
	_ = x[CallingConvHiPE-11]
	_ = x[CallingConvWebKitJS-12]
	_ = x[CallingConvAnyReg-13]
	_ = x[CallingConvPreserveMost-14]
	_ = x[CallingConvPreserveAll-15]
	_ = x[CallingConvSwift-16]
	_ = x[CallingConvCXXFastTLS-17]
	_ = x[CallingConvX86StdCall-64]
	_ = x[CallingConvX86FastCall-65]
	_ = x[CallingConvARM_APCS-66]
	_ = x[CallingConvARM_AAPCS-67]
	_ = x[CallingConvARM_AAPCS_VFP-68]
	_ = x[CallingConvMSP430Intr-69]
	_ = x[CallingConvX86ThisCall-70]
	_ = x[CallingConvPTXKernel-71]
	_ = x[CallingConvPTXDevice-72]
	_ = x[CallingConvSPIRFunc-75]
	_ = x[CallingConvSPIRKernel-76]
	_ = x[CallingConvIntelOCL_BI-77]
	_ = x[CallingConvX86_64SysV-78]
	_ = x[CallingConvWin64-79]
	_ = x[CallingConvX86VectorCall-80]
	_ = x[CallingConvHHVM-81]
	_ = x[CallingConvHHVM_C-82]
	_ = x[CallingConvX86Intr-83]
	_ = x[CallingConvAVRIntr-84]
	_ = x[CallingConvAVRSignal-85]
	_ = x[CallingConvAVRBuiltin-86]
	_ = x[CallingConvAMDGPU_VS-87]
	_ = x[CallingConvAMDGPU_GS-88]
	_ = x[CallingConvAMDGPU_PS-89]
	_ = x[CallingConvAMDGPU_CS-90]
	_ = x[CallingConvAMDGPUKernel-91]
	_ = x[CallingConvX86RegCall-92]
	_ = x[CallingConvAMDGPU_HS-93]
	_ = x[CallingConvMSP430Builtin-94]
	_ = x[CallingConvAMDGPU_LS-95]
	_ = x[CallingConvAMDGPU_ES-96]
	_ = x[CallingConvAArch64VectorCall-97]
}

const (
	_CallingConv_name_0 = "noneccc"
	_CallingConv_name_1 = "fastcccoldccghccccc 11webkit_jsccanyregccpreserve_mostccpreserve_allccswiftcccxx_fast_tlscc"
	_CallingConv_name_2 = "x86_stdcallccx86_fastcallccarm_apcsccarm_aapcsccarm_aapcs_vfpccmsp430_intrccx86_thiscallccptx_kernelptx_device"
	_CallingConv_name_3 = "spir_funcspir_kernelintel_ocl_biccx86_64_sysvccwin64ccx86_vectorcallcchhvmcchhvm_cccx86_intrccavr_intrccavr_signalcccc 86amdgpu_vsamdgpu_gsamdgpu_psamdgpu_csamdgpu_kernelx86_regcallccamdgpu_hscc 94amdgpu_lsamdgpu_esaarch64_vector_pcs"
)

var (
	_CallingConv_index_0 = [...]uint8{0, 4, 7}
	_CallingConv_index_1 = [...]uint8{0, 6, 12, 17, 22, 33, 41, 56, 70, 77, 91}
	_CallingConv_index_2 = [...]uint8{0, 13, 27, 37, 48, 63, 76, 90, 100, 110}
	_CallingConv_index_3 = [...]uint8{0, 9, 20, 34, 47, 54, 70, 76, 84, 94, 104, 116, 121, 130, 139, 148, 157, 170, 183, 192, 197, 206, 215, 233}
)

func (i CallingConv) String() string {
	switch {
	case i <= 1:
		return _CallingConv_name_0[_CallingConv_index_0[i]:_CallingConv_index_0[i+1]]
	case 8 <= i && i <= 17:
		i -= 8
		return _CallingConv_name_1[_CallingConv_index_1[i]:_CallingConv_index_1[i+1]]
	case 64 <= i && i <= 72:
		i -= 64
		return _CallingConv_name_2[_CallingConv_index_2[i]:_CallingConv_index_2[i+1]]
	case 75 <= i && i <= 97:
		i -= 75
		return _CallingConv_name_3[_CallingConv_index_3[i]:_CallingConv_index_3[i+1]]
	default:
		return "CallingConv(" + strconv.FormatInt(int64(i), 10) + ")"
	}
}
