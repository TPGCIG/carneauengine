import * as React from "react";
import { Check, ChevronsUpDown } from "lucide-react";
import { useVirtualizer } from "@tanstack/react-virtual";
import { clsx, type ClassValue } from "clsx";
import { twMerge } from "tailwind-merge";
import { E164Number } from "libphonenumber-js/core";

// ----------------------------------------------------------------------------
// UTILS & UI COMPONENTS
// ----------------------------------------------------------------------------

function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}

const Button = React.forwardRef<HTMLButtonElement, ButtonProps>(
  ({ className, variant = "default", ...props }, ref) => {
    return (
      <button
        ref={ref}
        className={cn(
          "inline-flex items-center justify-center whitespace-nowrap rounded-md text-sm font-medium ring-offset-background transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:pointer-events-none disabled:opacity-50 h-10 px-4 py-2",
          variant === "outline" && "border border-input bg-transparent hover:bg-accent hover:text-accent-foreground",
          className
        )}
        {...props}
      />
    );
  }
);
Button.displayName = "Button";

interface ButtonProps extends React.ButtonHTMLAttributes<HTMLButtonElement> {
  variant?: "outline" | "default"; // or any variants you want
}

const Input = React.forwardRef<HTMLInputElement, React.InputHTMLAttributes<HTMLInputElement>>
  (({ className, type, ...props }, ref) => (
    <input
      type={type}
      className={cn(
        "flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background file:border-0 file:bg-transparent file:text-sm file:font-medium placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50",
        className
      )}
      ref={ref}
      {...props}
    />
));
Input.displayName = "Input";

import {
  Command,
  CommandEmpty,
  CommandGroup,
  CommandInput,
  CommandItem,
  CommandList,
} from "cmdk";

import * as PopoverPrimitive from "@radix-ui/react-popover";

const Popover = PopoverPrimitive.Root;
const PopoverTrigger = PopoverPrimitive.Trigger;
const PopoverContent = React.forwardRef<
  React.ElementRef<typeof PopoverPrimitive.Content>,
  React.ComponentPropsWithoutRef<typeof PopoverPrimitive.Content>
>(({ className, align = "center", sideOffset = 4, ...props }, ref) => (
  <PopoverPrimitive.Portal>
    <PopoverPrimitive.Content
      ref={ref}
      align={align}
      sideOffset={sideOffset}
      className={cn(
        "z-50 w-72 rounded-none border bg-white p-4 text-popover-foreground shadow-md outline-none data-[state=open]:animate-in data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=open]:fade-in-0 data-[state=closed]:zoom-out-95 data-[state=open]:zoom-in-95 data-[side=bottom]:slide-in-from-top-2 data-[side=left]:slide-in-from-right-2 data-[side=right]:slide-in-from-left-2 data-[side=top]:slide-in-from-bottom-2",
        className
      )}
      {...props}
    />
  </PopoverPrimitive.Portal>
));
PopoverContent.displayName = PopoverPrimitive.Content.displayName;

// ----------------------------------------------------------------------------
// DATA: Country List (Embedded to avoid external dependencies)
// ----------------------------------------------------------------------------

interface CountryData {
  value: string; // ISO 2 code (e.g., 'US')
  label: string; // Name (e.g., 'United States')
  dial: string;  // Dial code (e.g., '1')
}

const COUNTRIES: CountryData[] = [
  { value: "US", label: "United States", dial: "1" },
  { value: "GB", label: "United Kingdom", dial: "44" },
  { value: "CA", label: "Canada", dial: "1" },
  { value: "AU", label: "Australia", dial: "61" },
  { value: "DE", label: "Germany", dial: "49" },
  { value: "FR", label: "France", dial: "33" },
  { value: "IN", label: "India", dial: "91" },
  { value: "JP", label: "Japan", dial: "81" },
  { value: "CN", label: "China", dial: "86" },
  { value: "BR", label: "Brazil", dial: "55" },
  { value: "MX", label: "Mexico", dial: "52" },
  { value: "IT", label: "Italy", dial: "39" },
  { value: "ES", label: "Spain", dial: "34" },
  { value: "RU", label: "Russia", dial: "7" },
  { value: "ZA", label: "South Africa", dial: "27" },
  { value: "KR", label: "South Korea", dial: "82" },
  { value: "NL", label: "Netherlands", dial: "31" },
  { value: "SE", label: "Sweden", dial: "46" },
  { value: "CH", label: "Switzerland", dial: "41" },
  { value: "TR", label: "Turkey", dial: "90" },
  { value: "ID", label: "Indonesia", dial: "62" },
  { value: "SA", label: "Saudi Arabia", dial: "966" },
  { value: "AR", label: "Argentina", dial: "54" },
  { value: "NZ", label: "New Zealand", dial: "64" },
  { value: "IE", label: "Ireland", dial: "353" },
  { value: "SG", label: "Singapore", dial: "65" },
  { value: "UA", label: "Ukraine", dial: "380" },
  { value: "EG", label: "Egypt", dial: "20" },
  { value: "TH", label: "Thailand", dial: "66" },
  { value: "VN", label: "Vietnam", dial: "84" },
  { value: "MY", label: "Malaysia", dial: "60" },
  { value: "PH", label: "Philippines", dial: "63" },
  { value: "PL", label: "Poland", dial: "48" },
  { value: "PK", label: "Pakistan", dial: "92" },
  { value: "NG", label: "Nigeria", dial: "234" },
  { value: "BD", label: "Bangladesh", dial: "880" },
  // Add more as needed. This list covers major regions.
];

// Helper to get emoji flag
function getFlagEmoji(countryCode: string) {
  const codePoints = countryCode
    .toUpperCase()
    .split('')
    .map(char =>  127397 + char.charCodeAt(0));
  return String.fromCodePoint(...codePoints);
}

// ----------------------------------------------------------------------------
// MAIN COMPONENT
// ----------------------------------------------------------------------------

export interface PhoneInputProps {
  value?: E164Number | string;
  onChange?: (value: E164Number | undefined) => void;
  defaultCountry?: string;
  className?: string;
}
const PhoneInput = React.forwardRef<HTMLInputElement, PhoneInputProps>(
  ({ className, onChange, value, defaultCountry = "US", ...props }, ref) => {
    // State to hold the selected country. 
    // We try to infer it from the value if possible, otherwise use default.
    const [selectedCountry, setSelectedCountry] = React.useState<CountryData>(
      COUNTRIES.find(c => c.value === defaultCountry) || COUNTRIES[0]
    );

    const [phoneNumber, setPhoneNumber] = React.useState(value || "");

    // Update internal state if prop changes
    React.useEffect(() => {
      if (value !== undefined) {
        setPhoneNumber(value);
      }
    }, [value]);

    const handleCountrySelect = (country: CountryData) => {
      setSelectedCountry(country);
      // Focus the input after selection
      // (optional, implementation dependent)
    };

    const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
      const newValue = e.target.value as E164Number;
      // Basic validation: allow only numbers, spaces, dash, parenthesis, plus
      if (/^[0-9+\-\s()]*$/.test(newValue)) {
        setPhoneNumber(newValue);
        if (onChange) onChange(newValue);
      }
    };

    return (
      <div className={cn("flex", className)}>
        <CountrySelect
          value={selectedCountry}
          onChange={handleCountrySelect}
          options={COUNTRIES}
        />
        <Input
          className="rounded-none placeholder:text-muted-foreground placeholder:opacity-50"
          type="tel"
          value={phoneNumber}
          onChange={handleInputChange}
          placeholder="123 456 789"
          {...props}
        />
      </div>
    );
  }
);
PhoneInput.displayName = "PhoneInput";


// ----------------------------------------------------------------------------
// COUNTRY SELECT COMPONENT
// ----------------------------------------------------------------------------

type CountrySelectProps = {
  disabled?: boolean;
  value: CountryData;
  onChange: (value: CountryData) => void;
  options: CountryData[];
};

const CountrySelect = ({
  disabled,
  value,
  onChange,
  options,
}: CountrySelectProps) => {
  const [open, setOpen] = React.useState(false);

  const handleSelect = React.useCallback(
    (country: CountryData) => {
      onChange(country);
      setOpen(false);
    },
    [onChange]
  );

  return (
    <Popover open={open} onOpenChange={setOpen}>
      <PopoverTrigger asChild>
        <Button
          type="button"
          variant="outline"
          className={cn("flex gap-1 rounded-noneroun px-3 border-r-0 focus:z-10 bg-transparent text-foreground border-input hover:bg-accent hover:text-accent-foreground w-[100px] justify-between")}
          disabled={disabled}
        >
          <div className="flex items-center gap-2 truncate">
             <span className="text-lg leading-none">{getFlagEmoji(value.value)}</span>
             <span className="text-xs text-muted-foreground">+{value.dial}</span>
          </div>
          <ChevronsUpDown
            className={cn(
              "h-4 w-4 opacity-50 shrink-0",
              disabled ? "hidden" : "opacity-100"
            )}
          />
        </Button>
      </PopoverTrigger>
      <PopoverContent className="w-[300px] p-0" align="start">
        <CountrySearchList 
            options={options} 
            selectedCountry={value} 
            onSelect={handleSelect} 
        />
      </PopoverContent>
    </Popover>
  );
};

// ----------------------------------------------------------------------------
// VIRTUALIZED SEARCH LIST
// ----------------------------------------------------------------------------

const CountrySearchList = ({
    options,
    selectedCountry,
    onSelect
}: {
    options: CountryData[],
    selectedCountry: CountryData,
    onSelect: (country: CountryData) => void
}) => {
    const [search, setSearch] = React.useState("");

    // Filter options based on search input
    const filteredOptions = React.useMemo(() => {
        if (!search) return options;
        const lowerSearch = search.toLowerCase();
        return options.filter((option) => 
            option.label.toLowerCase().includes(lowerSearch) || 
            option.value.toLowerCase().includes(lowerSearch) ||
            option.dial.includes(lowerSearch)
        );
    }, [options, search]);

    const parentRef = React.useRef<HTMLDivElement>(null);

    // Initialize virtualizer
    const rowVirtualizer = useVirtualizer({
        count: filteredOptions.length,
        getScrollElement: () => parentRef.current,
        estimateSize: () => 35, 
        overscan: 5,
    });

    return (
        <Command shouldFilter={false} className="h-full w-full overflow-hidden">
            <div className="flex items-center border-b px-3">
                <CommandInput 
                    placeholder="Search country..." 
                    value={search}
                    onValueChange={setSearch}
                    className="flex h-11 w-full rounded-none bg-transparent py-3 text-sm outline-none placeholder:text-muted-foreground disabled:cursor-not-allowed disabled:opacity-50"
                />
            </div>
            
            <CommandList className="max-h-[300px] overflow-y-auto overflow-x-hidden" ref={parentRef}>
                {filteredOptions.length === 0 && (
                    <CommandEmpty className="py-6 text-center text-sm">No country found.</CommandEmpty>
                )}

                <div
                    style={{
                        height: `${rowVirtualizer.getTotalSize()}px`,
                        width: '100%',
                        position: 'relative',
                    }}
                >
                    {rowVirtualizer.getVirtualItems().map((virtualRow) => {
                        const option = filteredOptions[virtualRow.index];
                        return (
                            <CommandItem
                              key={option.value}
                              value={option.label}
                              onSelect={() => onSelect(option)} 
                              onClick={() => onSelect(option)}
                              className="absolute left-0 top-0 w-full cursor-default select-none items-center rounded-none px-2 py-1.5 text-sm outline-none hover:bg-accent hover:text-accent-foreground flex gap-2"
                              style={{
                                height: `${virtualRow.size}px`,
                                transform: `translateY(${virtualRow.start}px)`,
                              }}
                            >
                              {/* Flag container with fixed width to prevent clipping */}
                              <span className="flex h-6 w-8 items-center justify-center text-lg">
                                {getFlagEmoji(option.value)}
                              </span>

                              <span className="flex-1 truncate">{option.label}</span>

                              <span className="text-foreground/50 text-xs">
                                +{option.dial}
                              </span>

                              {option.value === selectedCountry.value && (
                                <Check className="ml-auto h-4 w-4 opacity-100" />
                              )}
                            </CommandItem>

                        );
                    })}
                </div>
            </CommandList>
        </Command>
    );
};

export { PhoneInput };