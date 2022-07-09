section .data
b: 
    .long    4
c: 
    .long    3
d: 
    .long    5
a: 
    .long    2

result:
        .zero   4
.LC0:
        .string "%d"
fmt:
        .quad   .LC0
section .text

main:
	push    rbp
	mov     rbp, rsp
	sub     rsp, 16
    mov    ebx, DWORD PTR a[rip]
    mov    ebx, DWORD PTR b[rip]
    cdq
    idiv    ebx
    mov    DWORD PTR t0[rip], eax
    mov    ebx, DWORD PTR c[rip]
    mov    ebx, DWORD PTR d[rip]
    imul   eax, ebx
    mov    DWORD PTR t1[rip], eax
    mov    eax, DWORD PTR t1[rip]
    mov    eax, DWORD PTR t0[rip]
    add    eax, ebx
    mov    DWORD PTR result[rip], eax

	mov     edx, DWORD PTR result[rip]
	mov     rax, QWORD PTR fmt[rip]
	mov     esi, edx
	mov     rdi, rax
	mov     eax, 0
	call    printf
	mov     eax, 0
	leave
	ret
