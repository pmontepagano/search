package ar.com.montepagano.search.examples.creditcardpayments;

import java.io.FileInputStream;
import java.io.IOException;
import java.util.Scanner;
import io.grpc.ManagedChannel;
import io.grpc.ManagedChannelBuilder;
import java.util.Scanner;
import java.io.IOException;
import java.nio.file.Files;
import java.nio.file.Paths;
import java.io.FileInputStream;
import com.google.protobuf.ByteString;
import ar.com.montepagano.search.middleware.v1.PrivateMiddlewareServiceGrpc;
import ar.com.montepagano.search.middleware.v1.Middleware.RegisterChannelRequest;
import ar.com.montepagano.search.middleware.v1.Middleware.RegisterChannelResponse;
import ar.com.montepagano.search.contracts.v1.Contracts.GlobalContract;

// Press Shift twice to open the Search Everywhere dialog and type `show whitespaces`,
// then press Enter. You can now see whitespace characters in your code.
public class Main {
    public static void main(String[] args) {
        System.out.println( "Hello World!" );
        Scanner scanner = new Scanner(System.in);

        // List of book titles
        String[] bookTitles = {
                "The Great Gatsby",
                "To Kill a Mockingbird",
                "1984",
                "Pride and Prejudice",
                "Animal Farm",
                "The Catcher in the Rye",
                "Lord of the Flies",
                "Brave New World",
                "The Hobbit",
                "The Fellowship of the Ring"
        };

        // Prompt user to select book titles
        System.out.println("Please select at least one book to purchase:");
        for (int i = 0; i < bookTitles.length; i++) {
            System.out.println((i + 1) + ". " + bookTitles[i]);
        }

        // Get user input for selected book titles
        System.out.print("Enter the number(s) of the book(s) you want to purchase (comma-separated): ");
        String selectedBooksInput = scanner.nextLine();
        String[] selectedBooksInputArray = selectedBooksInput.split(",");
        int[] selectedBooks = new int[selectedBooksInputArray.length];
        for (int i = 0; i < selectedBooksInputArray.length; i++) {
            selectedBooks[i] = Integer.parseInt(selectedBooksInputArray[i].trim());
        }

        // Prompt user for shipping address
        System.out.print("Enter your shipping address: ");
        String shippingAddress = scanner.nextLine();

        // Print purchase information
        System.out.println("Purchase Information:");
        System.out.println("Selected Books:");
        for (int i = 0; i < selectedBooks.length; i++) {
            System.out.println("- " + bookTitles[selectedBooks[i] - 1]);
        }
        System.out.println("Shipping Address: " + shippingAddress);

        // Load file contract.fsa into a ByteString variable
        ByteString contractBytes = null;
        try {
            contractBytes = ByteString.readFrom(new FileInputStream("./contract.fsa"));
        } catch (IOException e) {
            e.printStackTrace();
        }
        GlobalContract contract = GlobalContract.newBuilder().setContract(contractBytes).build();

        // Register channel with middleware
        ManagedChannel channel = ManagedChannelBuilder.forTarget("localhost:1234").usePlaintext().build();
        PrivateMiddlewareServiceGrpc.PrivateMiddlewareServiceBlockingStub stub = PrivateMiddlewareServiceGrpc.newBlockingStub(channel);
        RegisterChannelRequest request = RegisterChannelRequest.newBuilder().setRequirementsContract(contract).build();
        RegisterChannelResponse response = stub.registerChannel(request);
        channel.shutdown();
    }
}