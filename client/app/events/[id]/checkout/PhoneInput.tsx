import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";

interface PhoneInputProps {
    areaCode: number;
    phoneNumber: number;
    onChange: (value: { areaCode: number; phoneNumber: number }) => void;
}

export default function PhoneInput({ areaCode, phoneNumber, onChange }: PhoneInputProps) {
  return (
    <div className="flex gap-2 items-center">
      <div>
        <Label htmlFor="areaCode">Area Code</Label>
        <Input
          id="areaCode"
          type="text"
          value={areaCode}
          maxLength={3}
          className="w-20"
          onChange={(e) => onChange({ areaCode: Number(e.target.value), phoneNumber })}
        />
      </div>

      <div>
        <Label htmlFor="phoneNumber">Number</Label>
        <Input
          id="phoneNumber"
          type="text"
          value={phoneNumber}
          className="flex-1"
          onChange={(e) => onChange({ areaCode, phoneNumber: Number(e.target.value) })}
        />
      </div>
    </div>
  );
}