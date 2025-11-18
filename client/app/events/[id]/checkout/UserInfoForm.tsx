"use client";
import { Card, CardContent, CardFooter } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { useState, useEffect } from "react";
import { ShoppingCart } from 'lucide-react';
import { Separator } from "@/components/ui/separator";
import  PhoneInput from "./PhoneInput"




import {
  Field,
  FieldDescription,
  FieldGroup,
  FieldLabel,
  FieldLegend,
  FieldSet,
} from "@/components/ui/field"
import { Input } from "@/components/ui/input"


export function UserInfoForm() {
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
            <Input id="street" type="text" placeholder="123 Main St" />
          </Field>
          <Field>
            <FieldLabel htmlFor="street">Street Address</FieldLabel>
            <PhoneInput areaCode={} phoneNumber={} onChange={} />
          </Field>
        </FieldGroup>
      </FieldSet>
    </div>
  )
}
