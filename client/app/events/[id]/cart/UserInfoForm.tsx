"use client";
import { Card, CardContent, CardFooter } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { useState, useEffect } from "react";
import { ShoppingCart } from 'lucide-react';
import { Separator } from "@/components/ui/separator";
import "react-phone-number-input/style.css"
import { PhoneInput } from "./PhoneInput"



import {
  Field,
  FieldDescription,
  FieldGroup,
  FieldLabel,
  FieldLegend,
  FieldSet,
} from "@/components/ui/field"
import { ScrollArea } from "@/components/ui/scroll-area";
import { Input } from "@/components/ui/input"
import { E164Number } from "libphonenumber-js/core";


export function UserInfoForm() {
    const [phone, setPhone] = useState<E164Number | undefined>()


    useEffect(()=>{console.log(phone)}, [phone])

    return (
      <div className="w-full max-w-md space-y-6">
        <FieldSet>
          <FieldLegend>Personal Information</FieldLegend>
          <FieldDescription>
            The tickets will be mailed to your email address.
          </FieldDescription>
          <FieldGroup>
            <Field>
              <FieldLabel htmlFor="street">Email Address</FieldLabel>
              <Input id="street" type="text" placeholder="Enter your email" className="outline-black placeholder:text-muted-foreground placeholder:opacity-50" />
            </Field>
            <Field>
              <FieldLabel htmlFor="street">Street Address</FieldLabel>
              <PhoneInput
                className="outline outline-offset-2 [outline-color:theme(colors.black)] rounded-xs"
                defaultCountry="AU"
                  value={phone}
                  onChange={setPhone}
                  placeholder="Enter your phone number"
                />
            </Field>
          </FieldGroup>
        </FieldSet>
      </div>
    )
  }
