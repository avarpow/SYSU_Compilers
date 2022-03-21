#include<stdio.h>
int main()
{
    int s=0,a,n,t;
    printf("input a n\n");
    scanf("%d%d", &a, &n);
    t=a;
    while(n>0)
    {
        s+=t;
        a=a*10;
        t+=a;
        n--;
        float b = 1.23e5;
        a=s+t;
    }
    printf("a+aa+...=%d\n",s);
    return 0;
}