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
          <FieldLegend><h2>Personal Information</h2></FieldLegend>
          <FieldDescription>
            The tickets will be mailed to your email address.
          </FieldDescription>
          <FieldGroup>
            <Field>
              <FieldLabel htmlFor="street" ><h3>Email Address:</h3></FieldLabel>
              <Input id="street" type="text" placeholder="user@email.com" className="bg-white rounded-none outline-none focus:outline-none focus:ring-0 shadow-none appearance-none placeholder:text-gray-500" />
            </Field>
            <Field>
              <FieldLabel htmlFor="street"><h3>Phone Number:</h3></FieldLabel>
              <PhoneInput
                className="rounded-none outline-none focus:outline-none focus:ring-0 shadow-none focus-within:shadow-none"
                defaultCountry="AU"
                  value={phone}
                  onChange={setPhone}
                  placeholder="412345678"
                />
            </Field>
          </FieldGroup>
        </FieldSet>
      </div>
    )
  }
