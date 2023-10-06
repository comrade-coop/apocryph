#![cfg_attr(not(feature = "std"), no_std, no_main)]

#[ink::contract]
mod payment {

    use scale::Compact;

    #[ink(storage)]
    pub struct Payment {
        value: Balance,
    }

    impl Payment {
        #[ink(constructor)]
        pub fn new(init_value: Compact<Balance>) -> Self {
            Self { value: init_value.into() }
        }

        #[ink(constructor)]
        pub fn default() -> Self {
            Self::new(Compact(Default::default()))
        }

        #[ink(message)]
        pub fn claim(&mut self, amount: Compact<Balance>) {
            self.value = self.value - amount.0;
        }

        #[ink(message)]
        pub fn get(&self) -> Compact<Balance> {
            Compact(self.value)
        }
    }

    #[cfg(test)]
    mod tests {
        use super::*;

        #[ink::test]
        fn default_works() {
            let payment = Payment::default();
            assert_eq!(payment.get(), 0);
        }
    }


    /*
    #[cfg(all(test, feature = "e2e-tests"))]
    mod e2e_tests {
        /// Imports all the definitions from the outer scope so we can use them here.
        use super::*;

        /// A helper function used for calling contract messages.
        use ink_e2e::build_message;

        /// The End-to-End test `Result` type.
        type E2EResult<T> = std::result::Result<T, Box<dyn std::error::Error>>;

        /// We test that we can upload and instantiate the contract using its default constructor.
        #[ink_e2e::test]
        async fn default_works(mut client: ink_e2e::Client<C, E>) -> E2EResult<()> {
            // Given
            let constructor = PaymentRef::default();

            // When
            let contract_account_id = client
                .instantiate("payment", &ink_e2e::alice(), constructor, 0, None)
                .await
                .expect("instantiate failed")
                .account_id;

            // Then
            let get = build_message::<PaymentRef>(contract_account_id.clone())
                .call(|payment| payment.get());
            let get_result = client.call_dry_run(&ink_e2e::alice(), &get, 0, None).await;
            assert!(matches!(get_result.return_value(), false));

            Ok(())
        }

        /// We test that we can read and write a value from the on-chain contract contract.
        #[ink_e2e::test]
        async fn it_works(mut client: ink_e2e::Client<C, E>) -> E2EResult<()> {
            // Given
            let constructor = PaymentRef::new(false);
            let contract_account_id = client
                .instantiate("payment", &ink_e2e::bob(), constructor, 0, None)
                .await
                .expect("instantiate failed")
                .account_id;

            let get = build_message::<PaymentRef>(contract_account_id.clone())
                .call(|payment| payment.get());
            let get_result = client.call_dry_run(&ink_e2e::bob(), &get, 0, None).await;
            assert!(matches!(get_result.return_value(), false));

            // When
            let flip = build_message::<PaymentRef>(contract_account_id.clone())
                .call(|payment| payment.flip());
            let _flip_result = client
                .call(&ink_e2e::bob(), flip, 0, None)
                .await
                .expect("flip failed");

            // Then
            let get = build_message::<PaymentRef>(contract_account_id.clone())
                .call(|payment| payment.get());
            let get_result = client.call_dry_run(&ink_e2e::bob(), &get, 0, None).await;
            assert!(matches!(get_result.return_value(), true));

            Ok(())
        }
    }*/
}
