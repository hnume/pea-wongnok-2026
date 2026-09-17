import  { NextRequest, NextResponse }  from "next/server"; 

export async function GET(req: NextRequest) {
        return NextResponse.json({ message: "hi I'm Oat"  });
} 
export async function  POST(req: NextRequest) {
   
const body = await  req.json();
   
    return  NextResponse.json({ received: body });
}