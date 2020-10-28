
export interface Address extends Extendable {
    e_scalar?: ScalarEnum;
    line1?: string;
    pin?: CustomScaler;
}

export type CustomScaler = number;

export interface Extendable {
    ext_property?: string;
}

export interface Person {
    address?: Array<Address>;
    age: number;
    friends: Array<string>;
    gender?: Person.GenderEnum;
    name?: string;
    prop1: boolean;
    prop2?: string;
    prop3?: number;
}

export namespace Person {
    export type GenderEnum = 'male' | 'female' | 'other' ;
    export const GenderEnum = {
        Male: 'male' as GenderEnum,
        Female: 'female' as GenderEnum,
        Other: 'other' as GenderEnum
    };
}

export type ScalarEnum = 'value1' | 'value2' ;
export const ScalarEnum = {
    Value1: 'value1' as ScalarEnum,
    Value2: 'value2' as ScalarEnum
};